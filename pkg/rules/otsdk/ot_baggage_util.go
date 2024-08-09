//go:build ignore

package baggage

type BaggageContainer struct {
	baggage interface{}
}

func (bc *BaggageContainer) TakeSnapShot() interface{} {
	return &BaggageContainer{bc.baggage}
}

func GetBaggageFromGLS() *Baggage {
	gls := GetBaggageContainerFromGLS()
	if gls == nil {
		return nil
	}
	p := gls.(*BaggageContainer).baggage
	if p != nil {
		return p.(*Baggage)
	} else {
		return nil
	}
}

func SetBaggageToGLS(baggage *Baggage) {
	SetBaggageContainerToGLS(&BaggageContainer{baggage})
}

func ClearBaggageInGLS() {
	SetBaggageToGLS(nil)
}

func DeleteBaggageMemberInGLS(key string) bool {
	if bInternal := GetBaggageFromGLS(); bInternal != nil {
		b := bInternal.DeleteMember(key)
		SetBaggageToGLS(&b)
		return true
	}
	return false
}
