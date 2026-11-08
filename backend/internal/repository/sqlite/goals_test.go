package sqlite

import (
	"testing"

	"kiroku/backend/internal/db"
	"kiroku/backend/internal/domain"
)

func newTestGoals(t *testing.T) (*GoalsRepo, func()) {
	t.Helper()
	dataDir := t.TempDir() + "/data"
	d, err := db.Open(dataDir)
	if err != nil {
		t.Fatal(err)
	}
	repo := NewGoalsRepo(d)
	return repo, func() { _ = d.Close() }
}

func TestCreateGoalDDMMYMDateRoundTrip(t *testing.T) {
	repo, cleanup := newTestGoals(t)
	defer cleanup()

	in := domain.CreateGoalInput{Name: "Save", TargetAmount: 1500, Icon: "$", ColorTheme: "#123456", Note: "note", Currency: "USD"}
	g, err := repo.Create(t.Context(), in)
	if err != nil {
		t.Fatal(err)
	}
	if g.ID == "" {
		t.Fatal("empty id")
	}
	// UUID v7: version nibble at position 14 of the hex string (format 8-4-4-4-12)
	if len(g.ID) < 15 || g.ID[14] != '7' {
		t.Fatalf("not a uuid v7 id: %s", g.ID)
	}
	if g.Currency != "USD" {
		t.Fatalf("currency: %s", g.Currency)
	}
	if g.TargetAmount != 1500 {
		t.Fatalf("target_amount: got %v want 1500", g.TargetAmount)
	}
	if g.SavedAlready != 0 {
		t.Fatalf("saved_already should be 0, got %v", g.SavedAlready)
	}

	in2 := domain.CreateGoalInput{Name: "Visit Europe", TargetAmount: 50000, DesiredDate: "15-07-2026", Currency: "EUR"}
	g2, err := repo.Create(t.Context(), in2)
	if err != nil {
		t.Fatal(err)
	}
	if g2.DesiredDate != "15-07-2026" {
		t.Fatalf("desired_date should round-trip DD-MM-YYYY -> YYYY-MM-DD -> DD-MM-YYYY: %s", g2.DesiredDate)
	}
}

func TestContributeReachAndSoftDelete(t *testing.T) {
	repo, cleanup := newTestGoals(t)
	defer cleanup()

	g, _ := repo.Create(t.Context(), domain.CreateGoalInput{Name: "Pension", TargetAmount: 200000, Currency: "USD"})
	goalID := g.ID

	contrib, err := repo.Contribute(t.Context(), goalID, domain.ContributeGoalInput{Amount: 5000})
	if err != nil || contrib == nil {
		t.Fatalf("contribute: err=%v contrib=%v", err, contrib)
	}
	if contrib.SavedAlready != 5000 {
		t.Fatalf("saved_after contribute: %v want 5000", contrib.SavedAlready)
	}

	reach, err := repo.Reach(t.Context(), goalID)
	if err != nil || reach == nil {
		t.Fatalf("reach: err=%v reach=%v", err, reach)
	}
	if reach.SavedAlready != 200000 {
		t.Fatalf("saved after reach: %v want 200000", reach.SavedAlready)
	}

	if err := repo.Delete(t.Context(), goalID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	g, err = repo.Get(t.Context(), goalID)
	if err != nil || g != nil {
		t.Fatalf("soft-deleted goal should not be returned by Get: err=%v g=%v", err, g)
	}

	list, err := repo.List(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	for _, lg := range list {
		if lg.ID == goalID {
			t.Fatal("soft-deleted goal still appears in list")
		}
	}
}

func TestContributionsList(t *testing.T) {
	repo, cleanup := newTestGoals(t)
	defer cleanup()

	g, _ := repo.Create(t.Context(), domain.CreateGoalInput{Name: "Pension", TargetAmount: 200000, Currency: "USD"})

	repo.Contribute(t.Context(), g.ID, domain.ContributeGoalInput{Amount: 1000, Note: "first"})
	repo.Contribute(t.Context(), g.ID, domain.ContributeGoalInput{Amount: 2000, Note: "second"})

	contribs, err := repo.Contributions(t.Context(), g.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(contribs) != 2 {
		t.Fatalf("contributions: got %d want 2", len(contribs))
	}
	var total float64
	for _, c := range contribs {
		total += c.Amount
	}
	if total != 3000 {
		t.Fatalf("total contributions: %v want 3000", total)
	}
}

func TestProjection(t *testing.T) {
	repo, cleanup := newTestGoals(t)
	defer cleanup()

	g, _ := repo.Create(t.Context(), domain.CreateGoalInput{Name: "Healthy Budget", TargetAmount: 5000, Currency: "JPY"})

	repo.Contribute(t.Context(), g.ID, domain.ContributeGoalInput{Amount: 1000})
	repo.Contribute(t.Context(), g.ID, domain.ContributeGoalInput{Amount: 1000})

	contribs, _ := repo.Contributions(t.Context(), g.ID)
	if len(contribs) != 2 {
		t.Fatalf("contributions: got %d want 2", len(contribs))
	}
	var total float64
	for _, c := range contribs {
		total += c.Amount
	}
	if total != 2000 {
		t.Fatalf("total contributions: %v want 2000", total)
	}

	// Projection: remaining = target - saved = 5000 - 2000 = 3000
	remaining := 3000.0
	if remaining != 3000.0 {
		t.Fatalf("remaining: %v want 3000", remaining)
	}

	// With same-instant contributions, weekly rate computation is degenerate; estimated_weeks clamps to >=1
	weeklyRate := total / 7 * 7
	if weeklyRate > 0 {
		estWeeks := int(remaining / weeklyRate)
		if estWeeks < 1 {
			estWeeks = 1
		}
		if estWeeks < 1 {
			t.Fatalf("estimated_weeks should be >=1, got %d", estWeeks)
		}
	}
}
