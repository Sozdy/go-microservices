package inventory

type partService struct {
	partRepository PartRepository
}

func NewPartService(partRepository PartRepository) *partService {
	return &partService{
		partRepository: partRepository,
	}
}
