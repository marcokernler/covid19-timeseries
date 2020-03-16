package internals

//
type CoVData struct {
	Items []CoVDataItem
}

//
type CoVDataItem struct {
	Province string
	Country  string
	Lat      string
	Lng      string
	Values   []CoVDataItemValue
}

//
type CoVDataItemValue struct {
	Date  string
	Cases string
	Deaths string
}

// get a single data-item by its province
func (c *CoVData) getByProvince(province string) CoVDataItem {
	//
	for _, covDataItem := range c.Items {
		//
		if covDataItem.Province == province {
			return covDataItem
		}
	}

	return CoVDataItem{}
}