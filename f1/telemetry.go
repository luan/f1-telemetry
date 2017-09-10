package f1

type TelemetryData struct {
	Time                 float32
	Laptime              float32
	Lapdistance          float32
	Totaldistance        float32
	X                    float32    // World space position
	Y                    float32    // World space position
	Z                    float32    // World space position
	Speed                float32    // Speed of car in m/s
	Xv                   float32    // Velocity in world space
	Yv                   float32    // Velocity in world space
	Zv                   float32    // Velocity in world space
	Xr                   float32    // World space right direction
	Yr                   float32    // World space right direction
	Zr                   float32    // World space right direction
	Xd                   float32    // World space forward direction
	Yd                   float32    // World space forward direction
	Zd                   float32    // World space forward direction
	SuspPos              [4]float32 // Note: All wheel arrays have the order:
	SuspVel              [4]float32 // RL, RR, FL, FR
	WheelSpeed           [4]float32
	Throttle             float32
	Steer                float32
	Brake                float32
	Clutch               float32
	Gear                 float32
	GforceLat            float32
	GforceLon            float32
	Lap                  float32
	Enginerate           float32
	SliProNativeSupport  float32    // SLI Pro support
	CarPosition          float32    // car race position
	KersLevel            float32    // kers energy left
	KersMaxLevel         float32    // kers maximum energy
	DRS                  float32    // 0 = off, 1 = on
	TractionControl      float32    // 0 (off) - 2 (high)
	AntiLockBrakes       float32    // 0 (off) - 1 (on)
	FuelInTank           float32    // current fuel mass
	FuelCapacity         float32    // fuel capacity
	InPits               float32    // 0 = none, 1 = pitting, 2 = in pit area
	Sector               float32    // 0 = sector1, 1 = sector2, 2 = sector3
	Sector1Time          float32    // time of sector1 (or 0)
	Sector2Time          float32    // time of sector2 (or 0)
	BrakesTemp           [4]float32 // brakes temperature (centigrade)
	TyresPressure        [4]float32 // tyres pressure PSI
	TeamInfo             float32    // team ID
	TotalLaps            float32    // total number of laps in this race
	TrackSize            float32    // track size meters
	LastLapTime          float32    // last lap time
	MaxRpm               float32    // cars max RPM, at which point the rev limiter will kick in
	IdleRpm              float32    // cars idle RPM
	MaxGears             float32    // maximum number of gears
	SessionType          float32    // 0 = unknown, 1 = practice, 2 = qualifying, 3 = race
	Drsallowed           float32    // 0 = not allowed, 1 = allowed, -1 = invalid / unknown
	TrackNumber          float32    // -1 for unknown, 0-21 for tracks
	Vehiclefiaflags      float32    // -1 = invalid/unknown, 0 = none, 1 = green, 2 = blue, 3 = yellow, 4 = red
	Era                  float32    // era, 2017 (modern) or 1980 (classic)
	EngineTemperature    float32    // engine temperature (centigrade)
	GforceVert           float32    // vertical g-force component
	AngVelX              float32    // angular velocity x-component
	AngVelY              float32    // angular velocity y-component
	AngVelZ              float32    // angular velocity z-component
	TyresTemperature     [4]byte    // tyres temperature (centigrade)
	TyresWear            [4]byte    // tyre wear percentage
	TyreCompound         byte       // compound of tyre – 0 = ultra soft, 1 = super soft, 2 = soft, 3 = medium, 4 = hard, 5 = inter, 6 = wet
	FrontBrakeBias       byte       // front brake bias (percentage)
	FuelMix              byte       // fuel mix - 0 = lean, 1 = standard, 2 = rich, 3 = max
	Currentlapinvalid    byte       // current lap invalid - 0 = valid, 1 = invalid
	TyresDamage          [4]byte    // tyre damage (percentage)
	FrontLeftWingDamage  byte       // front left wing damage (percentage)
	FrontRightWingDamage byte       // front right wing damage (percentage)
	RearWingDamage       byte       // rear wing damage (percentage)
	EngineDamage         byte       // engine damage (percentage)
	GearBoxDamage        byte       // gear box damage (percentage)
	ExhaustDamage        byte       // exhaust damage (percentage)
	PitLimiterStatus     byte       // pit limiter status – 0 = off, 1 = on
	PitSpeedLimit        byte       // pit speed limit in m/s
	SessionTimeLeft      float32    // NEW: time left in session in seconds
	RevLightsPercent     byte       // NEW: rev lights indicator (percentage)
	IsSpectating         byte       // NEW: whether the player is spectating
	SpectatorCarIndex    byte       // NEW: index of the car being spectated

	// Car data
	NumCars        byte        // number of cars in data
	PlayerCarIndex byte        // index of player's car in the array
	Cars           [20]CarData // data for all cars on track
}

type CarData struct {
	WorldPosition     [3]float32 // world co-ordinates of vehicle
	LastlapTime       float32
	CurrentlapTime    float32
	BestlapTime       float32
	Sector1Time       float32
	Sector2Time       float32
	LapDistance       float32
	DriverID          byte
	TeamID            byte
	CarPosition       byte // UPDATED: track positions of vehicle
	CurrentLapNum     byte
	TyreCompound      byte // compound of tyre – 0 = ultra soft, 1 = super soft, 2 = soft, 3 = medium, 4 = hard, 5 = inter, 6 = wet
	InPits            byte // 0 = none, 1 = pitting, 2 = in pit area
	Sector            byte // 0 = sector1, 1 = sector2, 2 = sector3
	Currentlapinvalid byte // current lap invalid - 0 = valid, 1 = invalid
	Penalties         byte // NEW: accumulated time penalties in seconds to be added
}
