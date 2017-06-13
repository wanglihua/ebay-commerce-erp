package order_logistics_service

type pending_inventory struct {
	stock_id   int64
	product_id int64
	qty        int64
}

type pending_inventory_temp_keeper struct {
	pending_inventory_list []pending_inventory
}

func new_pending_inventory_temp_keeper() *pending_inventory_temp_keeper {
	var instance = new(pending_inventory_temp_keeper)
	instance.pending_inventory_list = make([]pending_inventory, 0)

	return instance
}

func (c *pending_inventory_temp_keeper) add(stock_id int64, product_id int64, qty int64) {
	var pending_inventory_list = c.pending_inventory_list

	var pi_found pending_inventory

	var found = false
	var found_index = 0
	for index, pi := range pending_inventory_list {
		if pi.stock_id == stock_id && pi.product_id == product_id {
			pi_found = pi
			found = true
			found_index = index
			break
		}
	}

	if found == true {
		pi_found.qty = pi_found.qty + qty
		pending_inventory_list[found_index] = pi_found
	} else {
		pi_found.stock_id = stock_id
		pi_found.product_id = product_id
		pi_found.qty = qty
		pending_inventory_list = append(pending_inventory_list, pi_found)
	}
}

func (c *pending_inventory_temp_keeper) get(stock_id int64, product_id int64) int64 {
	var pending_inventory_list = c.pending_inventory_list
	for _, pi := range pending_inventory_list {
		if pi.stock_id == stock_id && pi.product_id == product_id {
			return pi.qty
		}
	}
	return 0
}
