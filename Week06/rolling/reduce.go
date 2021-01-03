package rolling

func Sum(iterator Iterator) float64 {
	result := 0.0
	for iterator.Next() {
		bucket := iterator.Bucket()
		for _, p := range bucket.Points {
			result += p
		}
	}
	return result
}

func Avg(iterator Iterator) float64 {
	result := 0.0
	count := 0.0
	for iterator.Next() {
		bucket := iterator.Bucket()
		for _, p := range bucket.Points {
			result += p
			count++
		}
	}
	return result / count
}

func Min(iterator Iterator) float64 {
	result := 0.0
	started := false
	for iterator.Next() {
		bucket := iterator.Bucket()
		for _, p := range bucket.Points {
			if !started {
				result = p
				started = true
				continue
			}
			if p < result {
				p = result
			}
		}
	}
	return result
}

func Max(iterator Iterator) float64 {
	result := 0.0
	started := false
	for iterator.Next() {
		bucket := iterator.Bucket()
		for _, p := range bucket.Points {
			if !started {
				result = p
				started = true
				continue
			}
			if p > result {
				p = result
			}
		}
	}
	return result
}

func Count(iterator Iterator) float64 {
	var result int64
	for iterator.Next() {
		bucket := iterator.Bucket()
		result += bucket.Count
	}
	return float64(result)
}
