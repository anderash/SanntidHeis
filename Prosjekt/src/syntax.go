type Elevator struct {
	IPADDR   string
	POSITION int
	/*
	   Etg.			Pos. nr.	ElevInfo.POSITION
	    1 ......... 0.........0
	  	  ......... 1.........0/1
		2 ......... 2.........1
		  ......... 3.........1/2
		3 ......... 4.........2
		  ......... 5.........2/3
		4 ......... 6.........3
	*/

	DIRECTION int
	/*
		opp = 1, ned = -1, stillestående = 0
	*/

	DESTINATION int
	/*
		1. etg = 0
		2. etg = 1
		3. etg = 2
		4. etg = 3
	*/

	ORDER_MATRIX [][]int
	/* 			   opp    	 ned    inne i heis			Settes til 1 ved en ordre
	1.etg	[[  0         0         0]
	2.etg 	 [  0         0         0]
	3.etg 	 [  0         0         0]
	4.etg	 [  0         0         0]]
	osv.
	*/
}

type ElevInfo struct {
	IPADDR     string
	F_NEW_INFO bool

	F_DEAD_ELEV   bool
	F_BUTTONPRESS bool
	F_ACK_ORDER bool

	POSITION    int
	DIRECTION   int
	DESTINATION int

	BUTTON_TYPE int
	BUTTONFLOOR int
}

type Output struct {
	OUTPUT_TYPE int
	/*
		LIGHT_OUTPUT = 0
		MOTOR_OUTPUT = 1
	*/

	LIGHT_TYPE int
	/*
		BUTTON_LAMP = 0
		FLOOR_INDICATOR = 1
		DOOR_LAMP = 2
	*/

	BUTTON_TYPE int
	/*
		BUTTON_CALL_UP = 0
	    BUTTON_CALL_DOWN = 1
	    BUTTON_COMMAND = 2
	    NO_BUTTON = -1
	*/

	FLOOR int

	VALUE int
	/*
		on = 1
		off = 0
	*/

	OUTPUT_DIRECTION int
	/*
		UP = 1
		STOP = 0
		DOWN = -1
	*/
}

type ButtonLight struct{
	BUTTON_TYPE int
	BUTTONFLOOR int
	VALUE int
}