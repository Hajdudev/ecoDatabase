package models

import (
	"time"
)

type Calendar struct {
	ServiceID string    `db:"service_id" json:"service_id"`
	Monday    bool      `db:"monday" json:"monday"`
	Tuesday   bool      `db:"tuesday" json:"tuesday"`
	Wednesday bool      `db:"wednesday" json:"wednesday"`
	Thursday  bool      `db:"thursday" json:"thursday"`
	Friday    bool      `db:"friday" json:"friday"`
	Saturday  bool      `db:"saturday" json:"saturday"`
	Sunday    bool      `db:"sunday" json:"sunday"`
	StartDate time.Time `db:"start_date" json:"start_date"`
	EndDate   time.Time `db:"end_date" json:"end_date"`
}

type CalendarDate struct {
	ServiceID     string    `db:"service_id" json:"service_id"`
	Date          time.Time `db:"date" json:"date"`
	ExceptionType int       `db:"exception_type" json:"exception_type"`
}

type Route struct {
	RouteID          string `db:"route_id" json:"route_id"`
	AgencyID         string `db:"agency_id" json:"agency_id"`
	RouteShortName   string `db:"route_short_name" json:"route_short_name"`
	RouteLongName    string `db:"route_long_name" json:"route_long_name"`
	RouteDescription string `db:"route_description" json:"route_description"`
	RouteType        int    `db:"route_type" json:"route_type"`
	RouteURL         string `db:"route_url" json:"route_url"`
	RouteColor       string `db:"route_color" json:"route_color"`
	RouteTextColor   string `db:"route_text_color" json:"route_text_color"`
	RouteSortOrder   int64  `db:"route_sort_order" json:"route_sort_order"`
}

type Shape struct {
	ShapeID           string  `db:"shape_id" json:"shape_id"`
	ShapePtLat        float64 `db:"shape_pt_lat" json:"shape_pt_lat"`
	ShapePtLon        float64 `db:"shape_pt_lon" json:"shape_pt_lon"`
	ShapePtSequence   int     `db:"shape_pt_sequence" json:"shape_pt_sequence"`
	ShapeDistTraveled float64 `db:"shape_dist_traveled" json:"shape_dist_traveled"`
}

type StopTime struct {
	TripID            string  `db:"trip_id" json:"trip_id"`
	ArrivalTime       string  `db:"arrival_time" json:"arrival_time"`
	DepartureTime     string  `db:"departure_time" json:"departure_time"`
	StopID            string  `db:"stop_id" json:"stop_id"`
	StopSequence      int     `db:"stop_sequence" json:"stop_sequence"`
	StopHeadsign      string  `db:"stop_headsign" json:"stop_headsign"`
	PickupType        int     `db:"pickup_type" json:"pickup_type"`
	DropOffType       int     `db:"drop_off_type" json:"drop_off_type"`
	ShapeDistTraveled float64 `db:"shape_dist_traveled" json:"shape_dist_traveled"`
	Timepoint         int     `db:"timepoint" json:"timepoint"`
}

type Stop struct {
	StopID             string  `db:"stop_id" json:"stop_id"`
	StopCode           string  `db:"stop_code" json:"stop_code"`
	StopName           string  `db:"stop_name" json:"stop_name"`
	StopDesc           string  `db:"stop_desc" json:"stop_desc"`
	StopLat            float64 `db:"stop_lat" json:"stop_lat"`
	StopLon            float64 `db:"stop_lon" json:"stop_lon"`
	ZoneID             string  `db:"zone_id" json:"zone_id"`
	StopURL            string  `db:"stop_url" json:"stop_url"`
	LocationType       int     `db:"location_type" json:"location_type"`
	ParentStation      string  `db:"parent_station" json:"parent_station"`
	StopTimezone       string  `db:"stop_timezone" json:"stop_timezone"`
	WheelchairBoarding int     `db:"wheelchair_boarding" json:"wheelchair_boarding"`
	LevelID            string  `db:"level_id" json:"level_id"`
	PlatformCode       string  `db:"platform_code" json:"platform_code"`
}

type Trip struct {
	RouteID              string `db:"route_id" json:"route_id"`
	ServiceID            string `db:"service_id" json:"service_id"`
	TripID               string `db:"trip_id" json:"trip_id"`
	TripHeadsign         string `db:"trip_headsign" json:"trip_headsign"`
	TripShortName        string `db:"trip_short_name" json:"trip_short_name"`
	DirectionID          int    `db:"direction_id" json:"direction_id"`
	BlockID              string `db:"block_id" json:"block_id"`
	ShapeID              string `db:"shape_id" json:"shape_id"`
	WheelchairAccessible int    `db:"wheelchair_accessible" json:"wheelchair_accessible"`
	BikesAllowed         int    `db:"bikes_allowed" json:"bikes_allowed"`
}

type User struct {
	ID          int64     `db:"id" json:"id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	Email       string    `db:"email" json:"email"`
	Name        string    `db:"name" json:"name"`
	Image       string    `db:"image" json:"image"`
	RecentRides []string  `db:"recent_rides" json:"recent_rides"`
}
