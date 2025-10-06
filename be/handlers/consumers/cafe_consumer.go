package consumers

import (
	"context"
	"fmt"
	"log"

	"github.com/RandySteven/CafeConnect/be/entities/messages"
	"github.com/RandySteven/CafeConnect/be/entities/payloads/responses"
	"github.com/RandySteven/CafeConnect/be/enums"
	cache_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/caches"
	repository_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/repositories"
	topics_interfaces "github.com/RandySteven/CafeConnect/be/interfaces/topics"
	"github.com/RandySteven/CafeConnect/be/utils"
)

type CafeConsumer struct {
	cafeTopic               topics_interfaces.CafeTopic
	cafeRepository          repository_interfaces.CafeRepository
	addressRepository       repository_interfaces.AddressRepository
	cafeFranchiseRepository repository_interfaces.CafeFranchiseRepository
	cafeCache               cache_interfaces.CafeCache
}

func (c *CafeConsumer) GetCafesByRadius(ctx context.Context) error {
	message, err := c.cafeTopic.ReadMessage(ctx)
	if err != nil {
		return err
	}
	cafeMessage := utils.ReadJSONObject[messages.CafeMessage](message)

	addresses, err := c.addressRepository.FindAddressBasedOnRadius(ctx, cafeMessage.Longitude, cafeMessage.Latitude, cafeMessage.Radius)
	if err != nil {
		return err
	}

	result := make([]*responses.ListCafeResponse, len(addresses))

	for _, address := range addresses {
		cafe, err := c.cafeRepository.FindByAddressId(ctx, address.ID)
		if err != nil {
			return err
		}
		cafeFranchise, err := c.cafeFranchiseRepository.FindByID(ctx, cafe.CafeFranchiseID)
		if err != nil {
			return err
		}

		result = append(result, &responses.ListCafeResponse{
			ID:        cafe.ID,
			Name:      cafeFranchise.Name,
			LogoURL:   cafeFranchise.LogoURL,
			Status:    utils.GetCafeOpenCloseStatus(cafe.OpenHour, cafe.CloseHour),
			OpenHour:  cafe.OpenHour,
			CloseHour: cafe.CloseHour,
			Address: struct {
				Address   string  `json:"address"`
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			}{Address: address.Address, Latitude: address.Latitude, Longitude: address.Longitude},
		})
	}

	key := fmt.Sprintf(enums.ListCafeRadiusKey, cafeMessage.Longitude, cafeMessage.Latitude, fmt.Sprintf("%d", cafeMessage.Radius))

	err = c.cafeCache.SetCafeRadiusListCache(ctx, key, result)
	if err != nil {
		log.Println(`failed to set cafe radius list cache`, err)
	}

	return nil
}

func newCafeConsumer(
	cafeTopic topics_interfaces.CafeTopic,
	cafeRepository repository_interfaces.CafeRepository,
	addressRepository repository_interfaces.AddressRepository,
	cafeFranchiseRepository repository_interfaces.CafeFranchiseRepository,
	cafeCache cache_interfaces.CafeCache,
) *CafeConsumer {
	return &CafeConsumer{
		cafeTopic:               cafeTopic,
		cafeRepository:          cafeRepository,
		addressRepository:       addressRepository,
		cafeFranchiseRepository: cafeFranchiseRepository,
		cafeCache:               cafeCache,
	}
}
