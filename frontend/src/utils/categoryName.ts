import { useI18n } from '@/locales'
import type { Category } from '@/types/category'

const NAME_TO_KEY: Record<string, string> = {
  // Japanese standard categories
  '食費': 'cat_food',
  '交通費': 'cat_transport',
  '住居費': 'cat_housing',
  '収入': 'cat_income',
  'その他': 'cat_other',
  '食材': 'cat_food_groceries',
  '外食': 'cat_food_dining',
  'カフェ': 'cat_food_cafe',
  'お菓子': 'cat_food_snacks',
  '飲料': 'cat_food_drinks',
  '電車': 'cat_transport_train',
  'バス': 'cat_transport_bus',
  'タクシー': 'cat_transport_taxi',
  'ガソリン': 'cat_transport_gas',
  '高速代': 'cat_transport_toll',
  '家賃': 'cat_housing_rent',
  '電気代': 'cat_housing_electric',
  'ガス代': 'cat_housing_gas',
  '水道代': 'cat_housing_water',
  '修繕費': 'cat_housing_maintenance',
  '給与': 'cat_income_salary',
  'ボーナス': 'cat_income_bonus',
  '副業': 'cat_income_side',
  '投資収益': 'cat_income_investment',
  'その他収入': 'cat_income_other',

  // English standard categories
  'Food': 'cat_food',
  'Transportation': 'cat_transport',
  'Housing': 'cat_housing',
  'Income': 'cat_income',
  'Other': 'cat_other',
  'Groceries': 'cat_food_groceries',
  'Dining': 'cat_food_dining',
  'Cafe': 'cat_food_cafe',
  'Snacks': 'cat_food_snacks',
  'Drinks': 'cat_food_drinks',
  'Train': 'cat_transport_train',
  'Bus': 'cat_transport_bus',
  'Taxi': 'cat_transport_taxi',
  'Gas': 'cat_transport_gas',
  'Toll': 'cat_transport_toll',
  'Rent': 'cat_housing_rent',
  'Electric': 'cat_housing_electric',
  'Water': 'cat_housing_water',
  'Maintenance': 'cat_housing_maintenance',
  'Salary': 'cat_income_salary',
  'Bonus': 'cat_income_bonus',
  'Side Hustle': 'cat_income_side',
  'Investment': 'cat_income_investment',
  'Other Income': 'cat_income_other',

  // Indonesian standard categories
  'Makanan': 'cat_food',
  'Transportasi': 'cat_transport',
  'Perumahan': 'cat_housing',
  'Pendapatan': 'cat_income',
  'Lainnya': 'cat_other',
  'Belanjaan': 'cat_food_groceries',
  'Makan Luar': 'cat_food_dining',
  'Kafe': 'cat_food_cafe',
  'Cemilan': 'cat_food_snacks',
  'Minuman': 'cat_food_drinks',
  'Kereta': 'cat_transport_train',
  'Taksi': 'cat_transport_taxi',
  'Bensin': 'cat_transport_gas',
  'Tol': 'cat_transport_toll',
  'Sewa': 'cat_housing_rent',
  'Listrik': 'cat_housing_electric',
  'Air': 'cat_housing_water',
  'Pemeliharaan': 'cat_housing_maintenance',
  'Gaji': 'cat_income_salary',
  'Sambilan': 'cat_income_side',
  'Investasi': 'cat_income_investment',
  'Pendapatan Lain': 'cat_income_other'
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
