import { useI18n } from '@/locales'
import type { Category } from '@/types/category'

const NAME_TO_KEY: Record<string, string> = {
  // English standard categories
  'Food & Drink': 'cat_food_drink',
  'Bar, Cafe': 'cat_bar_cafe',
  'Groceries': 'cat_groceries',
  'Restaurant, Fast-Food': 'cat_restaurant_fastfood',
  'Shopping': 'cat_shopping',
  'Clothes, & Shoes': 'cat_clothes_shoes',
  'Drug-Store, Chenist': 'cat_drug_store',
  'Electronics, Accessories': 'cat_electronics_accessories',
  'Free Time': 'cat_free_time',
  'Gifts & Joy': 'cat_gifts_joy',
  'Healthy & Beauty': 'cat_healthy_beauty',
  'Home, Garden': 'cat_home_garden',
  'Jewels & Accessories': 'cat_jewels_accessories',
  'Kids': 'cat_kids',
  'Pets, Animal': 'cat_pets_animal',
  'Investment': 'cat_investment',
  'Realty': 'cat_realty',
  'Savings': 'cat_savings',
  'Vehicles, Chattels': 'cat_vehicles_chattels',
  'Income': 'cat_income',
  'Rental Income': 'cat_rental_income',
  'Sale': 'cat_sale',
  'Wage, Invoices': 'cat_wage_invoices',
  'Dividends, Interest': 'cat_dividends_interest',
  'Checks, Coupons': 'cat_checks_coupons',
  'Gifts': 'cat_gifts_income',
  'Refund(Tax)': 'cat_refund_tax',
  'Child Support': 'cat_child_support_income',
  'Lottery, Gambling': 'cat_lottery_gambling_income',
  'Financial Expenses': 'cat_financial_expenses',
  'Taxes': 'cat_taxes',
  'Fines': 'cat_fines',
  'Loan, Interest': 'cat_loan_interest',
  'Insurances': 'cat_insurance_financial',
  'Advisory': 'cat_advisory',
  'Charges, Fees': 'cat_charges_fees',
  'Life & Entertainment': 'cat_life_entertainment',
  'TV, Streaming': 'cat_tv_streaming',
  'Active Sport & Fitness': 'cat_active_fitness',
  'Alcohol & Tobacco': 'cat_alcohol_tobacco',
  'Books, Audio, Subscriptions': 'cat_books_subscriptions',
  'Charity & Gifts': 'cat_charity_gifts',
  'Culture, Sport Events': 'cat_culture_sport_events',
  'Education, Development': 'cat_education_development',
  'Health Care, Doctor': 'cat_health_care_doctor',
  'Hobbies': 'cat_hobbies',
  'Holiday, Trips, Events': 'cat_holiday_trips_events',
  'Communication & PC': 'cat_communication_pc',
  'Internet': 'cat_internet',
  'Phone, Cell Phone': 'cat_phone_cell_phone',
  'Postal Services': 'cat_postal_services',
  'Software, Apps, Games': 'cat_software_apps_games',
  'Housing': 'cat_housing',
  'Energy, Utilities': 'cat_energy_utilities',
  'Maintenances, Repairs': 'cat_maintenances_repairs',
  'Mortgage': 'cat_mortgage',
  'Property Insurance': 'cat_property_insurance',
  'Rent': 'cat_rent',
  'Services': 'cat_services_housing',
  'Vehicle': 'cat_vehicle',
  'Fuel': 'cat_fuel',
  'Leasing': 'cat_leasing',
  'Parking': 'cat_parking',
  'Rentals': 'cat_rentals_vehicle',
  'Vehicle Insurance': 'cat_vehicle_insurance',
  'Vehicle Maintenance': 'cat_vehicle_maintenance',
  'Other': 'cat_other',

  // Japanese standard categories
  '飲食': 'cat_food_drink',
  'バー、カフェ': 'cat_bar_cafe',
  '食材': 'cat_groceries',
  'レストラン、ファストフード': 'cat_restaurant_fastfood',
  'ショッピング': 'cat_shopping',
  '衣類、靴': 'cat_clothes_shoes',
  'ドラッグストア': 'cat_drug_store',
  '電子機器、アクセサリ': 'cat_electronics_accessories',
  '余暇': 'cat_free_time',
  'ギフト、喜び': 'cat_gifts_joy',
  '健康、美容': 'cat_healthy_beauty',
  'ホーム、ガーデン': 'cat_home_garden',
  '宝石、アクセサリ': 'cat_jewels_accessories',
  'キッズ': 'cat_kids',
  'ペット、動物': 'cat_pets_animal',
  '投資': 'cat_investment',
  '不動産': 'cat_realty',
  '貯蓄': 'cat_savings',
  '車両、動産': 'cat_vehicles_chattels',
  '収入': 'cat_income',
  '賃貸収入': 'cat_rental_income',
  '販売': 'cat_sale',
  '賃金、請求書': 'cat_wage_invoices',
  '配当、利子': 'cat_dividends_interest',
  '小切手、クーポン': 'cat_checks_coupons',
  'ギフト': 'cat_gifts_income',
  '還付(税)': 'cat_refund_tax',
  '子供の養育費': 'cat_child_support_income',
  '宝くじ、ギャンブル': 'cat_lottery_gambling_income',
  '金融費用': 'cat_financial_expenses',
  '税金': 'cat_taxes',
  '罰金': 'cat_fines',
  'ローン、利息': 'cat_loan_interest',
  '保険': 'cat_insurance_financial',
  '顧問': 'cat_advisory',
  '手数料': 'cat_charges_fees',
  '生活、エンターテインメント': 'cat_life_entertainment',
  'テレビ、ストリーミング': 'cat_tv_streaming',
  'アクティブスポーツ、フィットネス': 'cat_active_fitness',
  'アルコール、タバコ': 'cat_alcohol_tobacco',
  '本、オーディオ、サブスクリプション': 'cat_books_subscriptions',
  '慈善、ギフト': 'cat_charity_gifts',
  '文化、スポーツイベント': 'cat_culture_sport_events',
  '教育、開発': 'cat_education_development',
  'ヘルスケア、医師': 'cat_health_care_doctor',
  '趣味': 'cat_hobbies',
  '休暇、旅行、イベント': 'cat_holiday_trips_events',
  '通信、PC': 'cat_communication_pc',
  'インターネット': 'cat_internet',
  '電話、携帯電話': 'cat_phone_cell_phone',
  '郵便サービス': 'cat_postal_services',
  'ソフトウェア、アプリ、ゲーム': 'cat_software_apps_games',
  '居住': 'cat_housing',
  'エネルギー、光熱費': 'cat_energy_utilities',
  'メンテナンス、修繕': 'cat_maintenances_repairs',
  '住宅ローン': 'cat_mortgage',
  '不動産保険': 'cat_property_insurance',
  '家賃': 'cat_rent',
  'サービス': 'cat_services_housing',
  '車両': 'cat_vehicle',
  '燃料': 'cat_fuel',
  'リース': 'cat_leasing',
  '駐車場': 'cat_parking',
  'レンタル': 'cat_rentals_vehicle',
  '車両保険': 'cat_vehicle_insurance',
  '車両メンテナンス': 'cat_vehicle_maintenance',
  'その他': 'cat_other',

  // Indonesian standard categories
  'Makanan & Minuman': 'cat_food_drink',
  'Bar, Kafe': 'cat_bar_cafe',
  'Belanjaan': 'cat_groceries',
  'Restoran, Makan Cepat': 'cat_restaurant_fastfood',
  'Belanja': 'cat_shopping',
  'Pakaian, & Sepatu': 'cat_clothes_shoes',
  'Toko Obat, Apotek': 'cat_drug_store',
  'Elektronik, Aksesoris': 'cat_electronics_accessories',
  'Waktu Luang': 'cat_free_time',
  'Hadiah & Kebahagiaan': 'cat_gifts_joy',
  'Sehat & Kecantikan': 'cat_healthy_beauty',
  'Rumah, Taman': 'cat_home_garden',
  'Perhiasan & Aksesoris': 'cat_jewels_accessories',
  'Anak-Anak': 'cat_kids',
  'Hewan Peliharaan': 'cat_pets_animal',
  'Investasi': 'cat_investment',
  'Properti': 'cat_realty',
  'Tabungan': 'cat_savings',
  'Kendaraan, Barang Berharga': 'cat_vehicles_chattels',
  'Pendapatan': 'cat_income',
  'Pemasukan Sewa': 'cat_rental_income',
  'Penjualan': 'cat_sale',
  'Gaji, Faktur': 'cat_wage_invoices',
  'Dividen, Bunga': 'cat_dividends_interest',
  'Cek, Kupon': 'cat_checks_coupons',
  'Hadiah': 'cat_gifts_income',
  'Pengembalian(Pajak)': 'cat_refund_tax',
  'Bantuan Anak': 'cat_child_support_income',
  'Lotre, Joker': 'cat_lottery_gambling_income',
  'Biaya Keuangan': 'cat_financial_expenses',
  'Pajak': 'cat_taxes',
  'Denda': 'cat_fines',
  'Pinjaman, Bunga': 'cat_loan_interest',
  'Asuransi': 'cat_insurance_financial',
  'Konsultasi': 'cat_advisory',
  'Biaya, Fee': 'cat_charges_fees',
  'Hidup & Hiburan': 'cat_life_entertainment',
  'TV Streaming': 'cat_tv_streaming',
  'Olahraga Aktif & Kebugaran': 'cat_active_fitness',
  'Alkohol & Tembakau': 'cat_alcohol_tobacco',
  'Buku, Audio, Langganan': 'cat_books_subscriptions',
  'Karuna & Hadiah': 'cat_charity_gifts',
  'Budaya, Acara Olahraga': 'cat_culture_sport_events',
  'Pendidikan, Pengembangan': 'cat_education_development',
  'Perawatan Kesehatan, Dokter': 'cat_health_care_doctor',
  'Hobi': 'cat_hobbies',
  'Liburan, Trip, Acara': 'cat_holiday_trips_events',
  'Komunikasi & PC': 'cat_communication_pc',
  'Jaringan Internet': 'cat_internet',
  'Telepon, HP': 'cat_phone_cell_phone',
  'Layanan Pos': 'cat_postal_services',
  'Software, Apps, Game': 'cat_software_apps_games',
  'Perumahan': 'cat_housing',
  'Energi, Utilitas': 'cat_energy_utilities',
  'Pemeliharaan, Perbaikan': 'cat_maintenances_repairs',
  'Kredit Rumah': 'cat_mortgage',
  'Asuransi Properti': 'cat_property_insurance',
  'Sewa': 'cat_rent',
  'Layanan': 'cat_services_housing',
  'Kendaraan': 'cat_vehicle',
  'Bahan Bakar': 'cat_fuel',
  'Sewa Kendaraan (Leasing)': 'cat_leasing',
  'Parkir': 'cat_parking',
  'Sewa Kendaraan': 'cat_rentals_vehicle',
  'Asuransi Kendaraan': 'cat_vehicle_insurance',
  'Pemeliharaan Kendaraan': 'cat_vehicle_maintenance',
  'Lainnya': 'cat_other',
}

/**
 * Returns helper functions to localize standard category names across languages.
 */
export function useCategoryName() {
  const { t } = useI18n()

  const translateName = (name: string): string => {
    if (!name) return ''
    if (name.includes(' / ')) {
      return name.split(' / ').map(translateName).join(' / ')
    }
    const mappedKey = NAME_TO_KEY[name]
    if (mappedKey) {
      const key = `standardCategories.${mappedKey}`
      const translated = t(key)
      if (translated !== key && translated !== '') {
        return translated
      }
    }
    return name
  }

  const translate = (category: Category | null | undefined): string => {
    if (!category) return ''

    // 1. Check if category.id matches a standard category id
    if (category.id) {
      const key = `standardCategories.${category.id}`
      const translated = t(key)
      if (translated !== key && translated !== '') {
        return translated
      }
    }

    // 2. Check path or name fallback
    if (category.path) {
      return translateName(category.path)
    }

    if (category.name) {
      return translateName(category.name)
    }

    return ''
  }

  return { translate, translateName }
}
