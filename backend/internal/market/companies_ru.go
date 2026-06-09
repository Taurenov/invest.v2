package market

// CompanyProfile — справочник эмитентов MOEX для обзоров и AI-прогноза.
type CompanyProfile struct {
	Symbol      string `json:"symbol"`
	Exchange    string `json:"exchange"`
	Name        string `json:"name"`
	Sector      string `json:"sector"`
	Description string `json:"description"`
	MarketCap   string `json:"market_cap,omitempty"`
}

var moexProfiles = map[string]CompanyProfile{
	"SBER": {
		Symbol: "SBER", Exchange: "MOEX", Name: "Сбербанк", Sector: "Финансы",
		Description: "Крупнейший банк России: розничное и корпоративное кредитование, экосистема (СберПрайм, маркетплейс, IT). Бенчмарк российского фондового рынка.",
		MarketCap: "топ-1 банковский сектор",
	},
	"GAZP": {
		Symbol: "GAZP", Exchange: "MOEX", Name: "Газпром", Sector: "Нефть и газ",
		Description: "Глобальный газовый холдинг: добыча, транспортировка и экспорт природного газа. Чувствителен к ценам на газ и геополитике.",
		MarketCap: "энергетический гигант",
	},
	"LKOH": {
		Symbol: "LKOH", Exchange: "MOEX", Name: "Лукойл", Sector: "Нефть и газ",
		Description: "Вертикально интегрированная нефтяная компания: разведка, добыча, переработка и розничные АЗС. Дивидендный профиль.",
		MarketCap: "топ нефтегаз",
	},
	"ROSN": {
		Symbol: "ROSN", Exchange: "MOEX", Name: "Роснефть", Sector: "Нефть и газ",
		Description: "Крупнейший нефтяной производитель РФ. Добыча на территории России и за рубежом, экспорт нефти.",
	},
	"GMKN": {
		Symbol: "GMKN", Exchange: "MOEX", Name: "Норникель", Sector: "Металлургия",
		Description: "Мировой лидер по производству никеля и палладия. Зависимость от цен на цветные металлы.",
	},
	"NVTK": {
		Symbol: "NVTK", Exchange: "MOEX", Name: "Новатэк", Sector: "Нефть и газ",
		Description: "Независимый производитель газа и СПГ. Проекты на Ямале и Арктике.",
	},
	"TATN": {
		Symbol: "TATN", Exchange: "MOEX", Name: "Татнефть", Sector: "Нефть и газ",
		Description: "Вертикально интегрированная нефтяная компания Татарстана. Добыча, переработка, нефтехимия.",
	},
	"YNDX": {
		Symbol: "YNDX", Exchange: "MOEX", Name: "Яндекс", Sector: "IT и технологии",
		Description: "Технологическая экосистема: поиск, такси, доставка, облако, маркетплейс. Рост зависит от digital-рынка.",
	},
	"VTBR": {
		Symbol: "VTBR", Exchange: "MOEX", Name: "ВТБ", Sector: "Финансы",
		Description: "Системно значимый банк: корпоративный и розничный бизнес, государственное участие.",
	},
	"MGNT": {
		Symbol: "MGNT", Exchange: "MOEX", Name: "Магнит", Sector: "Ритейл",
		Description: "Сеть продуктовых магазинов и дискаунтеров. Показатель потребительского спроса в регионах.",
	},
	"PLZL": {
		Symbol: "PLZL", Exchange: "MOEX", Name: "Полюс", Sector: "Металлургия",
		Description: "Крупнейший золотодобытчик России. Зависимость от цены на золото.",
	},
	"CHMF": {
		Symbol: "CHMF", Exchange: "MOEX", Name: "Северсталь", Sector: "Металлургия",
		Description: "Производитель стали и металлопродукции. Цикличный сектор, связан с промышленным спросом.",
	},
	"ALRS": {
		Symbol: "ALRS", Exchange: "MOEX", Name: "АЛРОСА", Sector: "Металлургия",
		Description: "Мировой лидер по добыче алмазов. Экспорт и внутренний ювелирный рынок.",
	},
	"MTSS": {
		Symbol: "MTSS", Exchange: "MOEX", Name: "МТС", Sector: "Телеком",
		Description: "Крупнейший мобильный оператор: связь, финтех, медиа и IT-сервисы.",
	},
	"MOEX": {
		Symbol: "MOEX", Exchange: "MOEX", Name: "Московская биржа", Sector: "Финансы",
		Description: "Оператор фондового, срочного и валютного рынков. Доходы от торговых оборотов и листинга.",
	},
}

func GetCompanyProfile(symbol, exchange string) (CompanyProfile, bool) {
	if exchange == "" {
		exchange = "MOEX"
	}
	p, ok := moexProfiles[symbol]
	if !ok {
		return CompanyProfile{}, false
	}
	return p, true
}

func ListMOEXProfiles() []CompanyProfile {
	order := []string{"SBER", "GAZP", "LKOH", "ROSN", "YNDX", "GMKN", "NVTK", "TATN", "VTBR", "MGNT", "PLZL", "CHMF", "ALRS", "MTSS", "MOEX"}
	out := make([]CompanyProfile, 0, len(order))
	for _, sym := range order {
		if p, ok := moexProfiles[sym]; ok {
			out = append(out, p)
		}
	}
	return out
}

func SummaryText(profile CompanyProfile, locale string) string {
	if locale == "en" {
		return profile.Name + " (" + profile.Symbol + "): " + profile.Description
	}
	return profile.Name + " (" + profile.Symbol + "): " + profile.Description
}
