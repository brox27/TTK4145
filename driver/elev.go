package driver

import (
	"../ConfigFile"
	"fmt"
)

var LAMP_CHANNEL_MATRIX = [ConfigFile.Num_floors][ConfigFile.Num_buttons]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4}}

var BUTTON_CHANNEL_MATRIX = [ConfigFile.Num_floors][ConfigFile.Num_buttons]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4}}

func InitElev() {
	if Io_init() == 0 {
		fmt.Println("Failed init")
	}

	for floor := 0; floor < ConfigFile.Num_floors; floor++ {
		for button := ConfigFile.BUTTON_ORDER_UP; button <= ConfigFile.Num_buttons; button++ {
			SetButtonLamp(button, floor, 0)
		}
	}

	SetStopLamp(0)
	SetDoorOpenLamp(0)
	SetFloorLight(0)
}

func SetMotorDirection(dir ConfigFile.Direction) {
	if dir == ConfigFile.NEUTRAL {
		Io_write_analog(MOTOR, 0)
	} else if dir == ConfigFile.UP {
		Io_clear_bit(MOTORDIR)
		Io_write_analog(MOTOR, ConfigFile.MOTOR_SPEED)
	} else if dir == ConfigFile.DOWN {
		Io_set_bit(MOTORDIR)
		Io_write_analog(MOTOR, ConfigFile.MOTOR_SPEED)
	}
}

func SetButtonLamp(button ConfigFile.ButtonType, floor int, value int) {
	if (ConfigFile.Num_floors > floor && floor >= 0) && (0 <= button && button < ConfigFile.Num_buttons) {
		if value == 1 {
			Io_set_bit(LAMP_CHANNEL_MATRIX[floor][button])
		} else {
			Io_clear_bit(LAMP_CHANNEL_MATRIX[floor][button])
		}
	}
}

func SetFloorLight(floor int) {
	if ConfigFile.Num_floors > floor && floor >= 0 {
		if floor&0x02 != 0 {
			Io_set_bit(LIGHT_FLOOR_IND1)
		} else {
			Io_clear_bit(LIGHT_FLOOR_IND1)
		}

		if floor&0x01 != 0 {
			Io_set_bit(LIGHT_FLOOR_IND2)
		} else {
			Io_clear_bit(LIGHT_FLOOR_IND2)
		}
	}
}

func SetDoorOpenLamp(value int) {
	if value == 1 {
		Io_set_bit(LIGHT_DOOR_OPEN)
	} else {
		Io_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func GetButtonSignal(floor int, button int) int {
	if (floor>=0 && floor < ConfigFile.Num_floors) && (button >= 0 && button < ConfigFile.Num_buttons){
		return Io_read_bit(BUTTON_CHANNEL_MATRIX[floor][button])
	}
	return 0
}

func GetFloorSensorSignal() int {
	if Io_read_bit(SENSOR_FLOOR1) != 0 {
		return 0
	}
	if Io_read_bit(SENSOR_FLOOR2) != 0 {
		return 1
	}
	if Io_read_bit(SENSOR_FLOOR3) != 0 {
		return 2
	}
	if Io_read_bit(SENSOR_FLOOR4 ) != 0 {
		return 3
	} else {
		return -1
	}
}

func SetStopLamp(value int) {
	if value != 0 {
		Io_set_bit(STOP)
	}
}
