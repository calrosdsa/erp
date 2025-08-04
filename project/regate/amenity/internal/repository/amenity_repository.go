package amenity_repo

type AmenityRepository interface {

}

type amenityRepository struct {

}

func NewAmenityRepository()AmenityRepository{
    return &amenityRepository{}
}

