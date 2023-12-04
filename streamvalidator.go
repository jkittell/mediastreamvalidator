package mediastreamvalidator

import (
	"encoding/json"
	"fmt"
	"github.com/jkittell/array"
	"github.com/jkittell/toolbox"
	"log"
	"time"
)

// content represents data about a video/audio stream.
type StreamValidator struct {
	Id         string         `json:"id"`
	URL        string         `json:"url"`
	Validation ValidationInfo `json:"validation_info"`
	Status     string         `json:"status"`
	StartTime  time.Time      `json:"start_time"`
	EndTime    time.Time      `json:"end_time"`
}

type ValidationInfo struct {
	IndependentSegments bool   `json:"independentSegments"`
	PlaylistKind        string `json:"playlistKind"`
	MimeType            string `json:"mimeType"`
	DataID              int    `json:"dataID"`
	GzipEncoded         bool   `json:"gzipEncoded"`
	ValidatorVersion    string `json:"validatorVersion"`
	ValidatorTimestamp  string `json:"validatorTimestamp"`
	Messages            []struct {
		ErrorComment          string `json:"errorComment"`
		ErrorDomain           string `json:"errorDomain"`
		ErrorStatusCode       int    `json:"errorStatusCode"`
		ErrorRequirementLevel int    `json:"errorRequirementLevel"`
		ErrorDetail           string `json:"errorDetail"`
		ErrorReferenceDataID  int    `json:"errorReferenceDataID,omitempty"`
	} `json:"messages"`
	SslContentDeliveredSecurely bool   `json:"sslContentDeliveredSecurely"`
	URL                         string `json:"url"`
	Variants                    []struct {
		MaxFrameRate           float64 `json:"maxFrameRate,omitempty"`
		VideoRangeKey          string  `json:"videoRangeKey,omitempty"`
		MimeType               string  `json:"mimeType"`
		ProcessedSegmentsCount int     `json:"processedSegmentsCount"`
		URL                    string  `json:"url"`
		HasDiscSequenceTag     bool    `json:"hasDiscSequenceTag"`
		Messages               []struct {
			ErrorComment          string `json:"errorComment"`
			ErrorDomain           string `json:"errorDomain"`
			ErrorStatusCode       int    `json:"errorStatusCode"`
			ErrorRequirementLevel int    `json:"errorRequirementLevel"`
			ErrorDetail           string `json:"errorDetail"`
		} `json:"messages"`
		MeanSegmentCount         int     `json:"meanSegmentCount"`
		MeasuredMeanBitrate      int     `json:"measuredMeanBitrate"`
		ParsedSegmentsCount      int     `json:"parsedSegmentsCount"`
		PlaylistMeanBitrate      int     `json:"playlistMeanBitrate,omitempty"`
		HasEndTag                bool    `json:"hasEndTag"`
		PlaylistCodecs           string  `json:"playlistCodecs,omitempty"`
		PlaylistResolutionHeight int     `json:"playlistResolutionHeight,omitempty"`
		MeanTotalDuration        float64 `json:"meanTotalDuration"`
		MeasuredMaxBitrate       int     `json:"measuredMaxBitrate"`
		IndependentSegments      bool    `json:"independentSegments"`
		AudioGroup               struct {
			PlaylistGroupID string `json:"playlistGroupID"`
			Renditions      []struct {
				URL          string `json:"url"`
				PersistentID int    `json:"persistentID"`
			} `json:"renditions"`
		} `json:"audioGroup,omitempty"`
		ClosedCaptionGroup struct {
			PlaylistGroupID string `json:"playlistGroupID"`
			Renditions      []struct {
				PersistentID int `json:"persistentID"`
			} `json:"renditions"`
		} `json:"closedCaptionGroup,omitempty"`
		IframeOnly                  bool   `json:"iframeOnly,omitempty"`
		PlaylistKind                string `json:"playlistKind"`
		PlaylistMaxBitrate          int    `json:"playlistMaxBitrate,omitempty"`
		SslContentDeliveredSecurely bool   `json:"sslContentDeliveredSecurely"`
		Discontinuities             []struct {
			Segments []struct {
				MediaSequence               int     `json:"mediaSequence"`
				SegmentDurationTag          float64 `json:"segmentDurationTag"`
				HasSampleAuxInfo            bool    `json:"hasSampleAuxInfo"`
				MimeType                    string  `json:"mimeType"`
				DataID                      int     `json:"dataID"`
				URL                         string  `json:"url"`
				SegmentByteRangeLength      int     `json:"segmentByteRangeLength"`
				VideoStartsWithIDR          bool    `json:"videoStartsWithIDR"`
				SslContentDeliveredSecurely bool    `json:"sslContentDeliveredSecurely"`
				SegmentByteRangeOffset      int     `json:"segmentByteRangeOffset"`
				Format                      string  `json:"format"`
				SegmentDateStamp            string  `json:"segmentDateStamp"`
				VideoFrameRate              float64 `json:"videoFrameRate"`
				StartTime                   float64 `json:"startTime"`
				DiscontinuityDomain         int     `json:"discontinuityDomain"`
			} `json:"segments"`
			Measurements struct {
				MeasuredMaxBitrate  float64 `json:"measuredMaxBitrate"`
				MeasuredMeanBitrate float64 `json:"measuredMeanBitrate"`
				MeasuredSegments    int     `json:"measuredSegments"`
			} `json:"measurements"`
			Tracks []struct {
				TrackVideoIDRStandardDeviation int     `json:"trackVideoIDRStandardDeviation,omitempty"`
				TrackID                        int     `json:"trackId"`
				TrackVideoWidth                int     `json:"trackVideoWidth,omitempty"`
				TrackVideoLevel                int     `json:"trackVideoLevel,omitempty"`
				TrackVideoTransferFunction     string  `json:"trackVideoTransferFunction,omitempty"`
				TrackVideoHeight               int     `json:"trackVideoHeight,omitempty"`
				TrackMediaType                 string  `json:"trackMediaType"`
				TrackVideoIDRInterval          float64 `json:"trackVideoIDRInterval,omitempty"`
				TrackVideoIsInterlaced         bool    `json:"trackVideoIsInterlaced,omitempty"`
				TrackVideoProfile              int     `json:"trackVideoProfile,omitempty"`
				TrackMediaSubType              string  `json:"trackMediaSubType"`
			} `json:"tracks"`
		} `json:"discontinuities"`
		PlaylistTargetDuration  int    `json:"playlistTargetDuration"`
		PlaylistResolutionWidth int    `json:"playlistResolutionWidth,omitempty"`
		GzipEncoded             bool   `json:"gzipEncoded"`
		DataID                  int    `json:"dataID"`
		PlaylistDefault         bool   `json:"playlistDefault,omitempty"`
		PlaylistName            string `json:"playlistName,omitempty"`
		PersistentID            int    `json:"persistentID,omitempty"`
		PlaylistLanguage        string `json:"playlistLanguage,omitempty"`
		PlaylistAutoselect      bool   `json:"playlistAutoselect,omitempty"`
		PlaylistMediaType       string `json:"playlistMediaType,omitempty"`
		PlaylistGroupID         string `json:"playlistGroupID,omitempty"`
		IsRendition             bool   `json:"isRendition,omitempty"`
		PlaylistChannels        string `json:"playlistChannels,omitempty"`
	} `json:"variants"`
	DataVersion float64 `json:"dataVersion"`
}

func Get(host, id string) StreamValidator {
	var info StreamValidator
	apiURL := fmt.Sprintf("http://%s:3001/contents/%s", host, id)
	status, res, err := toolbox.SendRequest(toolbox.GET, apiURL, "", nil)
	if err != nil {
		log.Println(err)
		return info
	}

	if status != 200 {
		log.Printf("get response code: %d", status)
	}

	err = json.Unmarshal(res, &info)
	if err != nil {
		log.Println(err)
	}
	return info
}

func GetAll(host string) *array.Array[StreamValidator] {
	all := array.New[StreamValidator]()
	apiURL := fmt.Sprintf("http://%s:3001/contents", host)
	status, res, err := toolbox.SendRequest(toolbox.GET, apiURL, "", nil)
	if err != nil {
		log.Println(err)
		return all
	}

	if status != 200 {
		log.Printf("get response code: %d", status)
	}

	fmt.Println(string(res))

	err = json.Unmarshal(res, &all)
	if err != nil {
		log.Println(err)
	}
	return all
}

func Post(host, url string) StreamValidator {
	var info StreamValidator
	apiURL := fmt.Sprintf("http://%s:3001/contents", host)

	data, _ := json.Marshal(map[string]string{"url": url})
	status, res, err := toolbox.SendRequest(toolbox.POST, apiURL, string(data), nil)
	if err != nil {
		log.Println(err)
		return info
	}

	if status != 201 {
		log.Printf("post response code: %d\n", status)
		return info
	}

	err = json.Unmarshal(res, &info)
	if err != nil {
		log.Println(err)
	}
	return info
}
