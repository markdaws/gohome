// NAA - NOAA's Astronomical Algorithms
package astrotime

// (JavaScript web page 
//  http://www.srrb.noaa.gov/highlights/sunrise/sunrise.html by
//  Chris Cornwall, Aaron Horiuchi and Chris Lehman)
// Ported to C++ by Pete Gray (petegray@ieee.org), July 2006
// Released as Open Source and can be used in any way, as long as the 
// above description remains in place.

import (
	"math"
	"time"
)

// Conversions
const (
	RadToDeg  = 180 / math.Pi
	DegToRad  = math.Pi / 180
	RadToGrad = 200 / math.Pi
	GradToDeg = math.Pi / 200
)

// More time constants
const (
	OneDay = time.Hour * 24
)

// CalcJD converts a time.Time object to a julian date
func CalcJD(t time.Time) float64 {
	y, m, d, hh, mm, ss, ms := t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1e6

	// Calc integer part (days)
	jday := (1461*(y+4800+(m-14)/12))/4 + (367*(m-2-12*((m-14)/12)))/12 - (3*((y+4900+(m-14)/12)/100))/4 + d - 32075

	// Calc floating point part (fraction of a day)
	jdatetime := float64(jday) + (float64(hh)-12.0)/24.0 + (float64(mm) / 1440.0) + (float64(ss) / 86400.0) + (float64(ms) / 86400000.0)

	// Adjust to UT
	_, zoneOffset := t.Zone()

	return jdatetime + float64(zoneOffset)/86400
}

// Name:    calcTimeJulianCent
// Type:    Function
// Purpose: convert Julian Day to centuries since J2000.0.
// Arguments:
//   jd : the Julian Day to convert
// Return value:
//   the T value corresponding to the Julian Day

func calcTimeJulianCent(t float64) float64 {
	return (t - 2451545.0) / 36525.0
}

// Name:    calcJDFromJulianCent
// Type:    Function
// Purpose: convert centuries since J2000.0 to Julian Day.
// Arguments:
//   t : number of Julian centuries since J2000.0
// Return value:
//   the Julian Day corresponding to the t value

func calcJDFromJulianCent(t float64) float64 {
	return t*36525.0 + 2451545.0
}

// Name:    calGeomMeanLongSun
// Type:    Function
// Purpose: calculate the Geometric Mean Longitude of the Sun
// Arguments:
//   t : number of Julian centuries since J2000.0
// Return value:
//   the Geometric Mean Longitude of the Sun in degrees

func calcGeomMeanLongSun(t float64) float64 {
	L0 := 280.46646 + t*(36000.76983+0.0003032*t)
	for L0 > 360.0 {
		L0 -= 360.0
	}
	for L0 < 0.0 {
		L0 += 360.0
	}
	return L0
}

// Name:    calcMeanObliquityOfEcliptic
// Type:    Function
// Purpose: calculate the mean obliquity of the ecliptic
// Arguments:
//   t : number of Julian centuries since J2000.0
// Return value:
//   mean obliquity in degrees

func calcMeanObliquityOfEcliptic(t float64) float64 {
	seconds := 21.448 - t*(46.8150+t*(0.00059-t*(0.001813)))
	return 23.0 + (26.0+(seconds/60.0))/60.0
}

// Name:    calcObliquityCorrection
// Type:    Function
// Purpose: calculate the corrected obliquity of the ecliptic
// Arguments:
//   t : number of Julian centuries since J2000.0
// Return value:
//   corrected obliquity in degrees

func calcObliquityCorrection(t float64) float64 {
	e0 := calcMeanObliquityOfEcliptic(t)
	omega := 125.04 - 1934.136*t
	return e0 + 0.00256*math.Cos(omega*DegToRad)
}

// Name:    calcEccentricityEarthOrbit
// Type:    Function
// Purpose: calculate the eccentricity of earth's orbit
// Arguments:
//   t : number of Julian centuries since J2000.0
// Return value:
//   the unitless eccentricity

func calcEccentricityEarthOrbit(t float64) float64 {
	return 0.016708634 - t*(0.000042037+0.0000001267*t)
}

// Name:    calGeomAnomalySun
// Type:    Function
// Purpose: calculate the Geometric Mean Anomaly of the Sun
// Arguments:
//   t : number of Julian centuries since J2000.0
// Return value:
//   the Geometric Mean Anomaly of the Sun in degrees

func calcGeomMeanAnomalySun(t float64) float64 {
	return 357.52911 + t*(35999.05029-0.0001537*t)
}

// Name:    calcEquationOfTime
// Type:    Function
// Purpose: calculate the difference between true solar time and mean
//		solar time
//Arguments:
//  t : number of Julian centuries since J2000.0
// Return value:
//   equation of time in minutes of time

func calcEquationOfTime(t float64) float64 {
	epsilon := calcObliquityCorrection(t)
	l0 := calcGeomMeanLongSun(t)
	e := calcEccentricityEarthOrbit(t)
	m := calcGeomMeanAnomalySun(t)

	y := math.Tan(DegToRad * epsilon / 2.0)
	y *= y

	sin2l0 := math.Sin(2.0 * DegToRad * l0)
	sinm := math.Sin(DegToRad * m)
	cos2l0 := math.Cos(2.0 * DegToRad * l0)
	sin4l0 := math.Sin(4.0 * DegToRad * l0)
	sin2m := math.Sin(2.0 * DegToRad * m)

	Etime := y*sin2l0 - 2.0*e*sinm + 4.0*e*y*sinm*cos2l0 - 0.5*y*y*sin4l0 - 1.25*e*e*sin2m

	return RadToDeg * Etime * 4.0
}

// Name:    calcSunEqOfCenter
// Type:    Function
// Purpose: calculate the equation of center for the sun
// Arguments:
//   t : number of Julian centuries since J2000.0
// Return value:
//   in degrees

func calcSunEqOfCenter(t float64) float64 {
	m := calcGeomMeanAnomalySun(t)
	mrad := DegToRad * m
	sinm := math.Sin(mrad)
	sin2m := math.Sin(mrad + mrad)
	sin3m := math.Sin(mrad + mrad + mrad)
	return sinm*(1.914602-t*(0.004817+0.000014*t)) + sin2m*(0.019993-0.000101*t) + sin3m*0.000289
}

// Name:    calcSunTrueLong
// Type:    Function
// Purpose: calculate the true longitude of the sun
// Arguments:
//   t : number of Julian centuries since J2000.0
// Return value:
//   sun's true longitude in degrees

func calcSunTrueLong(t float64) float64 {
	l0 := calcGeomMeanLongSun(t)
	c := calcSunEqOfCenter(t)
	return l0 + c
}

// Name:    calcSunApparentLong
// Type:    Function
// Purpose: calculate the apparent longitude of the sun
// Arguments:
//   t : number of Julian centuries since J2000.0
// Return value:
//   sun's apparent longitude in degrees

func calcSunApparentLong(t float64) float64 {
	o := calcSunTrueLong(t)
	omega := 125.04 - 1934.136*t
	return o - 0.00569 - 0.00478*math.Sin(DegToRad*omega)
}

// Name:    calcSunDeclination
// Type:    Function
// Purpose: calculate the declination of the sun
// Arguments:
//   t : number of Julian centuries since J2000.0
// Return value:
//   sun's declination in degrees

func calcSunDeclination(t float64) float64 {
	e := calcObliquityCorrection(t)
	lambda := calcSunApparentLong(t)
	sint := math.Sin(DegToRad*e) * math.Sin(DegToRad*lambda)
	return RadToDeg * math.Asin(sint)
}

// Name:    calcHourAngleSunrise
// Type:    Function
// Purpose: calculate the hour angle of the sun at sunrise for the
//			latitude
//Arguments:
//   lat : latitude of observer in degrees
//	solarDec : declination angle of sun in degrees
//Return value:
//   hour angle of sunrise in radians

func calcHourAngleSunrise(lat float64, solarDec float64) float64 {
	latRad := DegToRad * lat
	sdRad := DegToRad * solarDec
	return (math.Acos(math.Cos(DegToRad*90.833)/(math.Cos(latRad)*math.Cos(sdRad)) - math.Tan(latRad)*math.Tan(sdRad)))
}

// Name:    calcSolNoonUTC
// Type:    Function
// Purpose: calculate the Universal Coordinated Time (UTC) of solar
//		noon for the given day at the given location on earth
// Arguments:
//   t : number of Julian centuries since J2000.0
//   longitude : longitude of observer in degrees
// Return value:
//   time in minutes from zero Z

func calcSolNoonUTC(t float64, longitude float64) float64 {
	// First pass uses approximate solar noon to calculate eqtime
	tnoon := calcTimeJulianCent(calcJDFromJulianCent(t) + longitude/360.0)
	eqTime := calcEquationOfTime(tnoon)
	solNoonUTC := 720 + (longitude * 4) - eqTime
	newt := calcTimeJulianCent(calcJDFromJulianCent(t) - 0.5 + solNoonUTC/1440.0)
	eqTime = calcEquationOfTime(newt)
	return 720 + (longitude * 4) - eqTime
}

// Name:    calcSunriseUTC
// Type:    Function
// Purpose: calculate the Universal Coordinated Time (UTC) of sunrise
//			for the given day at the given location on earth
// Arguments:
//   JD  : julian day
//   latitude : latitude of observer in degrees
//   longitude : longitude of observer in degrees
// Return value:
//   time in minutes from zero Z

// Calculate the UTC sunrise for the given day at the given location
func calcSunriseUTC(jd float64, latitude float64, longitude float64) float64 {

	t := calcTimeJulianCent(jd)

	// *** Find the time of solar noon at the location, and use
	//     that declination. This is better than start of the
	//     Julian day

	noonmin := calcSolNoonUTC(t, longitude)
	tnoon := calcTimeJulianCent(jd + noonmin/1440.0)

	// *** First pass to approximate sunrise (using solar noon)

	eqTime := calcEquationOfTime(tnoon)
	solarDec := calcSunDeclination(tnoon)
	hourAngle := calcHourAngleSunrise(latitude, solarDec)

	delta := longitude - RadToDeg*hourAngle
	timeDiff := 4 * delta
	timeUTC := 720 + timeDiff - eqTime

	// *** Second pass includes fractional jday in gamma calc

	newt := calcTimeJulianCent(calcJDFromJulianCent(t) + timeUTC/1440.0)
	eqTime = calcEquationOfTime(newt)
	solarDec = calcSunDeclination(newt)
	hourAngle = calcHourAngleSunrise(latitude, solarDec)
	delta = longitude - RadToDeg*hourAngle
	timeDiff = 4 * delta
	timeUTC = 720 + timeDiff - eqTime
	return timeUTC
}

// CalcSunrise calculates the sunrise, in local time,  on the day t at the 
// location specified in longitude and latitude.
func CalcSunrise(t time.Time, latitude float64, longitude float64) time.Time {
	jd := CalcJD(t)
	sunriseUTC := time.Duration(math.Floor(calcSunriseUTC(jd, latitude, longitude)*60) * 1e9)
	loc, _ := time.LoadLocation("UTC")
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc).Add(sunriseUTC).In(t.Location())
}

// Name:    calcHourAngleSunset
// Type:    Function
// Purpose: calculate the hour angle of the sun at sunset for the
//			latitude
// Arguments:
//   lat : latitude of observer in degrees
//	solarDec : declination angle of sun in degrees
// Return value:
//   hour angle of sunset in radians

func calcHourAngleSunset(lat float64, solarDec float64) float64 {
	latRad := DegToRad * lat
	sdRad := DegToRad * solarDec

	HA := (math.Acos(math.Cos(DegToRad*90.833)/(math.Cos(latRad)*math.Cos(sdRad)) - math.Tan(latRad)*math.Tan(sdRad)))

	return -HA // in radians
}

// Name:    calcSunsetUTC
// Type:    Function
// Purpose: calculate the Universal Coordinated Time (UTC) of sunset
//			for the given day at the given location on earth
//Arguments:
//   JD  : julian day
//   latitude : latitude of observer in degrees
//   longitude : longitude of observer in degrees
// Return value:
//   time in minutes from zero Z

func calcSunsetUTC(jd float64, latitude float64, longitude float64) float64 {
	t := calcTimeJulianCent(jd)

	// *** Find the time of solar noon at the location, and use
	//     that declination. This is better than start of the
	//     Julian day

	noonmin := calcSolNoonUTC(t, longitude)
	tnoon := calcTimeJulianCent(jd + noonmin/1440.0)

	// First calculates sunrise and approx length of day

	eqTime := calcEquationOfTime(tnoon)
	solarDec := calcSunDeclination(tnoon)
	hourAngle := calcHourAngleSunset(latitude, solarDec)

	delta := longitude - RadToDeg*hourAngle
	timeDiff := 4 * delta
	timeUTC := 720 + timeDiff - eqTime

	// first pass used to include fractional day in gamma calc

	newt := calcTimeJulianCent(calcJDFromJulianCent(t) + timeUTC/1440.0)
	eqTime = calcEquationOfTime(newt)
	solarDec = calcSunDeclination(newt)
	hourAngle = calcHourAngleSunset(latitude, solarDec)

	delta = longitude - RadToDeg*hourAngle
	timeDiff = 4 * delta
	return 720 + timeDiff - eqTime
}

// CalcSunset calculates the sunset, in local time,  on the day t at the 
// location specified in longitude and latitude.
func CalcSunset(t time.Time, latitude float64, longitude float64) time.Time {
	jd := CalcJD(t)
	sunsetUTC := time.Duration(math.Floor(calcSunsetUTC(jd, latitude, longitude)*60) * 1e9)
	loc, _ := time.LoadLocation("UTC")
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc).Add(sunsetUTC).In(t.Location())
}

// NextSunrise returns date/time of the next sunrise after tAfter
func NextSunrise(tAfter time.Time, latitude float64, longitude float64) (tSunrise time.Time) {
	tSunrise = CalcSunrise(tAfter, latitude, longitude)
	if tAfter.After(tSunrise) {
		tSunrise = CalcSunrise(tAfter.Add(OneDay), latitude, longitude)
	}
	return
}

// NextSunset returns date/time of the next sunset after tAfter
func NextSunset(tAfter time.Time, latitude float64, longitude float64) (tSunset time.Time) {
	tSunset = CalcSunset(tAfter, latitude, longitude)
	if tAfter.After(tSunset) {
		tSunset = CalcSunset(tAfter.Add(OneDay), latitude, longitude)
	}
	return
}
