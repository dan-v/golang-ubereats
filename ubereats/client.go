package ubereats

import (
	"strings"
	"net/http"
	"encoding/json"
	"github.com/pkg/errors"
)

const (
	urlGetCsrfAndCookie = "https://www.ubereats.com/rtapi/locations/v2/predictions"
	urlStoreList        = "https://www.ubereats.com/rtapi/eats/v1/bootstrap-eater"
	urlStoreDetail      = "https://www.ubereats.com/rtapi/eats/v2/stores/"
)

type UberEatsClient struct {
	csrfToken string
	cookie    string
	latitude  string
	longitude string
}

func NewClient(latitude, longitude string) (*UberEatsClient, error) {
	u := &UberEatsClient{
		latitude:  latitude,
		longitude: longitude,
	}
	err := u.setup()
	return u, err
}

func (u *UberEatsClient) setup() error {
	// initialize request
	req, err := http.NewRequest("GET", urlGetCsrfAndCookie, nil)
	if err != nil {
		return errors.Wrap(err, "Failed to initialize request to " + urlGetCsrfAndCookie)
	}

	// set required headers
	req.Header.Set("Content-Type", "application/json")

	// make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "Failed to make request to " + urlGetCsrfAndCookie)
	}
	defer resp.Body.Close()

	// set csrf token and cookie
	u.csrfToken = resp.Header.Get("X-Csrf-Token")
	u.cookie = resp.Header.Get("Set-Cookie")

	return nil
}

func (u *UberEatsClient) GetStoreList() (Stores, error) {
	// initialize request
	payload := `{"targetLocation":{"latitude":` + u.latitude + `,"longitude":` + u.longitude + `}}`
	body := strings.NewReader(payload)
	req, err := http.NewRequest("POST", urlStoreList, body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize request to " + urlStoreList + " with payload " + payload)
	}

	// set required headers
	req.Header.Set("X-Csrf-Token", u.csrfToken)
	req.Header.Set("Cookie", u.cookie)
	req.Header.Set("Content-Type", "application/json")

	// make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to make request to " + urlStoreList + " with payload " + payload)
	}
	defer resp.Body.Close()

	// decode json to struct
	s := &StoreList{}
	if err := json.NewDecoder(resp.Body).Decode(s); err != nil {
		return nil, errors.Wrap(err, "Failed to decode GetStoreList response")
	}

	// check for empty
	if s.Marketplace.CityName == "" {
		return nil, errors.Wrap(err, "Failed to decode GetStoreList response")
	}

	return s.GetStores(), nil
}

func (u *UberEatsClient) GetStoreDetails(storeUUID string) (*StoreDetails, error) {
	// initialize request
	url :=  urlStoreDetail + storeUUID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize request to " + url)
	}

	// set required headers
	req.Header.Set("Cookie", u.cookie)

	// make request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to make request to " + url)
	}
	defer resp.Body.Close()

	// decode json to struct
	sd := &StoreDetails{}
	if err := json.NewDecoder(resp.Body).Decode(sd); err != nil {
		return nil, errors.Wrap(err, "Failed to decode GetStoreList response")
	}

	// check for empty
	if sd.Store.SubsectionsMap == nil {
		return nil, errors.Wrap(err, "Failed to decode GetStoreList response")
	}

	return sd, nil
}

type Store struct {
	UUID   string `json:"uuid"`
	Status string `json:"status"`
	Title  string `json:"title"`
	Categories []struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
	} `json:"categories,omitempty"`
	Tags []struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
	} `json:"tags,omitempty"`
	PriceBucket       string `json:"priceBucket,omitempty"`
	HeroImageURL      string `json:"heroImageUrl"`
	LargeHeroImageURL string `json:"largeHeroImageUrl,omitempty"`
	Endorsement struct {
		BackgroundColor struct {
			Alpha int `json:"alpha"`
			Color string `json:"color"`
		} `json:"backgroundColor"`
		IconColor struct {
			Alpha int `json:"alpha"`
			Color string `json:"color"`
		} `json:"iconColor"`
		IconURL string `json:"iconUrl"`
		Text    string `json:"text"`
		TextColor struct {
			Alpha int `json:"alpha"`
			Color string `json:"color"`
		} `json:"textColor"`
	} `json:"endorsement,omitempty"`
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Address struct {
			Address1   string `json:"address1"`
			AptOrSuite string `json:"aptOrSuite"`
			City       string `json:"city"`
			Country    string `json:"country"`
			PostalCode string `json:"postalCode"`
			Region     string `json:"region"`
			Title      string `json:"title"`
		} `json:"address"`
	} `json:"location"`
	RegionID            int `json:"regionId"`
	NotOrderableMessage string `json:"notOrderableMessage"`
	EtaRange struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"etaRange"`
	IsOrderable      bool `json:"isOrderable"`
	IsStoreMenuOpen  bool `json:"isStoreMenuOpen"`
	IsStoreVisible   bool `json:"isStoreVisible"`
	ClosedEtaMessage string `json:"closedEtaMessage"`
}

type Stores []Store

type StoreList struct {
	Experiments []interface{} `json:"experiments"`
	Marketplace struct {
		Timezone                      string `json:"timezone"`
		LearnMoreURL                  string `json:"learnMoreUrl"`
		PriceFormat                   string `json:"priceFormat"`
		CurrencyDecimalSeparator      string `json:"currencyDecimalSeparator"`
		CurrencyNumDigitsAfterDecimal int `json:"currencyNumDigitsAfterDecimal"`
		CurrencyCode                  string `json:"currencyCode"`
		AllowCredits                  bool `json:"allowCredits"`
		MyEats struct {
			Title        string `json:"title"`
			MyEatsStores []interface{} `json:"myEatsStores"`
			More struct {
				Tagline         string `json:"tagline"`
				DescriptionText string `json:"descriptionText"`
			} `json:"more"`
		} `json:"myEats"`
		Support struct {
			ContactPhone       string `json:"contactPhone"`
			ContactEmail       string `json:"contactEmail"`
			Title              string `json:"title"`
			EmailButtonText    string `json:"emailButtonText"`
			PhoneButtonText    string `json:"phoneButtonText"`
			ContactSupportText string `json:"contactSupportText"`
		} `json:"support"`
		Stores                           Stores `json:"stores"`
		IsInServiceArea                  bool `json:"isInServiceArea"`
		CityName                         string `json:"cityName"`
		MarketplaceCheckoutDeliveryTitle string `json:"marketplaceCheckoutDeliveryTitle"`
		Search struct {
			SuggestedSearches []struct {
				Name string `json:"name"`
				UUID string `json:"uuid"`
			} `json:"suggestedSearches"`
		} `json:"search"`
		Filters []struct {
			UUID string `json:"uuid"`
			Options []struct {
				UUID     string `json:"uuid"`
				Value    string `json:"value"`
				Selected bool `json:"selected"`
				Title    string `json:"title"`
			} `json:"options"`
			MinPermitted int `json:"minPermitted"`
			MaxPermitted int `json:"maxPermitted"`
			Type         string `json:"type"`
			Title        string `json:"title"`
		} `json:"filters"`
		Billboards []struct {
			UUID            string `json:"uuid"`
			Title           string `json:"title"`
			HeroImageURL    string `json:"heroImageUrl"`
			Type            string `json:"type"`
			MaxDisplayCount int `json:"maxDisplayCount"`
			StartTime       int `json:"startTime"`
			EndTime         int64 `json:"endTime"`
			Subtitle        string `json:"subtitle"`
			Link            string `json:"link"`
		} `json:"billboards"`
		DeliveryHoursInfos []struct {
			Date string `json:"date"`
			OpenHours []struct {
				StartTime      int `json:"startTime"`
				EndTime        int `json:"endTime"`
				DurationOffset int `json:"durationOffset"`
				IncrementStep  int `json:"incrementStep"`
			} `json:"openHours"`
		} `json:"deliveryHoursInfos"`
		Feed struct {
		} `json:"feed"`
	} `json:"marketplace"`
	Meta struct {
		TargetLocation struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"targetLocation"`
		Hashes struct {
			Instant string `json:"instant"`
			Stores  string `json:"stores"`
			Orders  string `json:"orders"`
			Client  string `json:"client"`
		} `json:"hashes"`
	} `json:"meta"`
	Orders []interface{} `json:"orders"`
	Tabs   []interface{} `json:"tabs"`
}

func (s *StoreList) GetStores() Stores {
	return s.Marketplace.Stores
}

type MenuItems map[string]struct {
	ImageURL        string `json:"imageUrl"`
	ItemDescription string `json:"itemDescription"`
	Price           int `json:"price"`
	Title           string `json:"title"`
	UUID            string `json:"uuid"`
	Options         []interface{} `json:"options"`
	Customizations []struct {
		UUID  string `json:"uuid"`
		Title string `json:"title"`
		Tags []struct {
			UUID string `json:"uuid"`
			Name string `json:"name"`
		} `json:"tags"`
		Options []struct {
			UUID  string `json:"uuid"`
			Title string `json:"title"`
			Price int `json:"price"`
			Tags []struct {
				UUID string `json:"uuid"`
				Name string `json:"name"`
			} `json:"tags"`
		} `json:"options"`
		MinPermitted int `json:"minPermitted"`
		MaxPermitted int `json:"maxPermitted"`
	} `json:"customizations"`
	AlcoholicItems int `json:"alcoholicItems"`
}

func (m MenuItems) String() string {
	output := ""
	for _, v := range m {
		output += v.Title + "\n"
	}
	return output
}

type StoreDetails struct {
	Meta struct {
		Hashes struct {
			Store string `json:"store"`
		} `json:"hashes"`
	} `json:"meta"`
	Store struct {
		UUID string `json:"uuid"`
		Location struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Address struct {
				Address1   string `json:"address1"`
				AptOrSuite string `json:"aptOrSuite"`
				City       string `json:"city"`
				Country    string `json:"country"`
				PostalCode string `json:"postalCode"`
				Region     string `json:"region"`
				Title      string `json:"title"`
			} `json:"address"`
		} `json:"location"`
		Categories []struct {
			UUID string `json:"uuid"`
			Name string `json:"name"`
		} `json:"categories"`
		ClosedEtaMessage string `json:"closedEtaMessage"`
		Endorsement struct {
			BackgroundColor struct {
				Alpha int `json:"alpha"`
				Color string `json:"color"`
			} `json:"backgroundColor"`
			IconColor struct {
				Alpha int `json:"alpha"`
				Color string `json:"color"`
			} `json:"iconColor"`
			IconURL string `json:"iconUrl"`
			Text    string `json:"text"`
			TextColor struct {
				Alpha int `json:"alpha"`
				Color string `json:"color"`
			} `json:"textColor"`
		} `json:"endorsement"`
		HeroImageURL        string `json:"heroImageUrl"`
		IsOrderable         bool `json:"isOrderable"`
		MenuItems           MenuItems `json:"itemsMap"`
		NotOrderableMessage string `json:"notOrderableMessage"`
		PriceBucket         string `json:"priceBucket"`
		Sections []struct {
			UUID     string `json:"uuid"`
			Title    string `json:"title"`
			Subtitle string `json:"subtitle"`
			IsTop    bool `json:"isTop"`
			IsOnSale bool `json:"isOnSale"`
			SubsectionGroups []struct {
				Title                  string `json:"title"`
				Type                   string `json:"type"`
				SubsectionDisplayOrder []string `json:"subsectionDisplayOrder"`
			} `json:"subsectionGroups"`
		} `json:"sections"`
		SubsectionsMap map[string]struct {
			UUID  string `json:"uuid"`
			Title string `json:"title"`
			DisplayItems []struct {
				UUID string `json:"uuid"`
				Type string `json:"type"`
			} `json:"displayItems"`
		} `json:"subsectionsMap"`
		Tags []struct {
			UUID string `json:"uuid"`
			Name string `json:"name"`
		} `json:"tags"`
		Title             string `json:"title"`
		RegionID          int `json:"regionId"`
		LargeHeroImageURL string `json:"largeHeroImageUrl"`
		Status            string `json:"status"`
		IsStoreVisible    bool `json:"isStoreVisible"`
		IsStoreMenuOpen   bool `json:"isStoreMenuOpen"`
		ShoppingCartItemsMap struct {
		} `json:"shoppingCartItemsMap"`
		DeliveryHoursInfos []struct {
			Date string `json:"date"`
			OpenHours []struct {
				StartTime      int `json:"startTime"`
				EndTime        int `json:"endTime"`
				DurationOffset int `json:"durationOffset"`
				IncrementStep  int `json:"incrementStep"`
			} `json:"openHours"`
		} `json:"deliveryHoursInfos"`
		SurgeInfo struct {
		} `json:"surgeInfo"`
		CanScheduleOrder bool `json:"canScheduleOrder"`
		SellsAlcohol     bool `json:"sellsAlcohol"`
	} `json:"store"`
}

func (s *StoreDetails) GetItems() MenuItems {
	return s.Store.MenuItems
}

func (s *StoreDetails) GetPopularItems() MenuItems {
	items := MenuItems{}
	if s.Store.SubsectionsMap != nil {
		for _, subsection := range s.Store.SubsectionsMap {
			if subsection.Title == "Most Popular" {
				for _, item := range subsection.DisplayItems {
					items[item.UUID] = s.Store.MenuItems[item.UUID]
				}
				break
			}
		}
	}
	return items
}