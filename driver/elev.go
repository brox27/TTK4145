package driver
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
#include "channels.h"
*/

import(
	"C"
	. "../globals"
)


func setMotorDirection(elev *Elevator, dir Direction){
	if (elev.dir == IDLE){
		Io_write_analog(MOTOR, 0);
    } else if (elev.dir == UP) {
        Io_clear_bit(MOTORDIR);
        Io_write_analog(MOTOR, MOTOR_SPEED);
    } else if (elev.dir == DOWN) {
        Io_set_bit(MOTORDIR);
        Io_write_analog(MOTOR, MOTOR_SPEED);
    }
}


func setButtonLamp(button ButtonType, floor int, value int){
	if ((N_FLOORS>=floor>0)&&(0<button<=N_BUTTONS)){
		if (value){
			Io_set_bit(LAMP_CHANNEL_MATRIX[floor][button])
		}
		else{
			Io_clear_bit(LAMP_CHANNEL_MATRIX[floor][button])
		}
	}
}


func setFloorLight(floor int){
	if (N_FLOORS>=floor>0){
		if (floor & 0x02) {
        Io_set_bit(LIGHT_FLOOR_IND1);
    	} 
    	else {
        Io_clear_bit(LIGHT_FLOOR_IND1);
    	}

    	if (floor & 0x01) {
        Io_set_bit(LIGHT_FLOOR_IND2);
    	} 
    	else {
        Io_clear_bit(LIGHT_FLOOR_IND2);
    	}    
	}
}


func setDoorOpenLamp(value int){
	if(value){
		Io_set_bit(LIGHT_DOOR_OPEN)
	}
	else{
		Io_clear_bit(LIGHT_DOOR_OPEN)
	}
}


func getButtonSignal(button ButtonType, floor int){
	if((N_FLOORS>=floor>0)&&(0<button<=N_BUTTONS)){
		return Io_read_bit(BUTTON_CHANNEL_MATRIX[floor][button])
	}
}


func getFloorSensorSignal(){
	if(Io_read_bit(SENSOR_FLOOR1)){
		return 1
	}
	if(Io_read_bit(SENSOR_FLOOR2)){
		return 2
	}
	if(Io_read_bit(SENSOR_FLOOR3)){
		return 3
	}
	if(Io_read_bit(SENSOR_FLOOR4)){
		return 4
	}
	else{ 
		return -1
	}
}
