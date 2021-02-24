package main

import (
	"fmt"
	"time"
	"math"
	"encoding/json"
)

type calculation_slice struct {
	Start_x float64
	Start_time time.Time
	Current_x float64
	R float64
	Iteration_counter int
	Tfs time.Duration
	Cycle []float64
	Cycle_list []float64
}

const rounding_error_limit float64 = 0.01
const cycle_accept_criteria int = 5
const cycle_maximum_size int = 10000


func get_next_x(f func (float64, float64) float64, current_x, r float64) float64 {
    return f(current_x, r)
}


func basic_generator(current_x, r float64) float64 {
    return current_x * r * (1 - current_x)
}


func is_it_equals(first_value, second_value float64) bool {
    
    value := math.Pow(math.Abs((math.Pow(first_value, 2.0) - math.Pow(second_value, 2.0))), 1.0/2.0)

    return value < rounding_error_limit
}

func is_equals_arrays(first_array, second_array []float64) bool {
    
    if len(first_array) != len(second_array) {
        return false
    }

    for i := 0; i <len(first_array); i++ {

        if !is_it_equals(first_array[i], second_array[i]) {
            return false
        }
    }

    return true
}

func get_json_from_calculation_slice(calculation_slice_map *calculation_slice) string{
	    b, err := json.Marshal(calculation_slice_map)
    	
    	if err != nil {
        	fmt.Println(err)
        	return "{}"
    	}
    return string(b)
}

func get_cycle_values(generator func (float64, float64) float64, calculation_slice_map *calculation_slice, start_x, r float64) []float64 {

    current_x := start_x
    calculation_slice_map.Start_x = current_x
    calculation_slice_map.R = r

    var cycle_list []float64
    cycle_list = append(cycle_list, current_x)

    iteration_counter := 0
    for {
        calculation_slice_map.Iteration_counter =  iteration_counter
        calculation_slice_map.Current_x = current_x
        calculation_slice_map.Cycle_list = cycle_list
        
        next_x := get_next_x(generator, current_x, r)
        cycle_list[iteration_counter] = next_x


        for i := 1;  i < len(cycle_list) / cycle_accept_criteria; i++ {
            calculation_slice_map.Tfs = time.Since(calculation_slice_map.Start_time)

            fmt.Println(get_json_from_calculation_slice(calculation_slice_map))

            cycle := cycle_list[len(cycle_list) - i:]

            place_to_cycle_search := cycle_list[len(cycle_list) - cycle_accept_criteria * i:]

            calculation_slice_map.Cycle = cycle
            
            cycle_attempt_count := 0
            
            for j := 0; j < cycle_accept_criteria; j++ {

                first_cycle_attempt := place_to_cycle_search[:i]
                place_to_cycle_search = place_to_cycle_search[i:]
                
                if is_equals_arrays(cycle, first_cycle_attempt) {
                    cycle_attempt_count += 1
                    
                    if cycle_attempt_count == cycle_accept_criteria {
                        return cycle
                    }
                }
            }
        }

        if len(cycle_list) > cycle_maximum_size * cycle_accept_criteria {
            return cycle_list
        }

        current_x = next_x
        cycle_list = append(cycle_list, current_x)

        iteration_counter += 1
    }
}

func main() {
	calculation_slice_map := &calculation_slice {
	
	}

	calculation_slice_map.Start_time = time.Now()

    for r := 0.0; r < 3.8; r += 0.01 {
        get_cycle_values(basic_generator, calculation_slice_map, 0.4, r)
    }

	fmt.Println("Thats all")
}
