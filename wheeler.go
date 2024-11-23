package main

type ColorSetting struct {
	Color   string `json:"color"`
	Enabled bool   `json:"enabled"`
}
type Entry struct {
	Text string `json:"text"`
}
type SpinnyWheel struct {
	AfterSpinSound             string         `json:"afterSpinSound"`
	AfterSpinSoundVolume       int            `json:"afterSpinSoundVolume"`
	AllowDuplicates            bool           `json:"allowDuplicates"`
	AnimateWinner              bool           `json:"animateWinner"`
	AutoRemoveWinner           bool           `json:"autoRemoveWinner"`
	CenterText                 string         `json:"centerText"`
	ColorSettings              []ColorSetting `json:"colorSettings"`
	CoverImageName             string         `json:"coverImageName"`
	CoverImageType             string         `json:"coverImageType"`
	CustomCoverImageDataURI    string         `json:"customCoverImageDataUri"`
	CustomPictureDataURI       string         `json:"customPictureDataUri"`
	CustomPictureName          string         `json:"customPictureName"`
	Description                string         `json:"description"`
	DisplayHideButton          bool           `json:"displayHideButton"`
	DisplayRemoveButton        bool           `json:"displayRemoveButton"`
	DisplayWinnerDialog        bool           `json:"displayWinnerDialog"`
	DrawOutlines               bool           `json:"drawOutlines"`
	DrawShadow                 bool           `json:"drawShadow"`
	DuringSpinSound            string         `json:"duringSpinSound"`
	DuringSpinSoundVolume      int            `json:"duringSpinSoundVolume"`
	Entries                    []Entry        `json:"entries"`
	GalleryPicture             string         `json:"galleryPicture"`
	HubSize                    string         `json:"hubSize"`
	IsAdvanced                 bool           `json:"isAdvanced"`
	LaunchConfetti             bool           `json:"launchConfetti"`
	MaxNames                   int            `json:"maxNames"`
	PageBackgroundColor        string         `json:"pageBackgroundColor"`
	PictureType                string         `json:"pictureType"`
	PlayClickWhenWinnerRemoved bool           `json:"playClickWhenWinnerRemoved"`
	ShowTitle                  bool           `json:"showTitle"`
	SlowSpin                   bool           `json:"slowSpin"`
	SpinTime                   int            `json:"spinTime"`
	Title                      string         `json:"title"`
	Type                       string         `json:"type"`
	WinnerMessage              string         `json:"winnerMessage"`
}

func GenerateWheel(entryList []string) SpinnyWheel {
	var entries []Entry
	for _, entry := range entryList {
		entries = append(entries, Entry{Text: entry})
	}
	wh := SpinnyWheel{
		AfterSpinSound:       "applause-sound-soft",
		AfterSpinSoundVolume: 50,
		AllowDuplicates:      false,
		AnimateWinner:        false,
		AutoRemoveWinner:     false,
		CenterText:           "",
		ColorSettings: []ColorSetting{
			{Color: "#3369E8", Enabled: true},
			{Color: "#D50F25", Enabled: true},
			{Color: "#EEB211", Enabled: true},
			{Color: "#009925", Enabled: true},
			{Color: "#000000", Enabled: false},
			{Color: "#000000", Enabled: false},
		},
		CoverImageName:             "",
		CoverImageType:             "",
		CustomCoverImageDataURI:    "",
		CustomPictureDataURI:       "",
		CustomPictureName:          "",
		Description:                "",
		DisplayHideButton:          true,
		DisplayRemoveButton:        true,
		DisplayWinnerDialog:        true,
		DrawOutlines:               false,
		DrawShadow:                 true,
		DuringSpinSound:            "ticking-sound",
		DuringSpinSoundVolume:      50,
		Entries:                    entries,
		GalleryPicture:             "/images/none.png",
		HubSize:                    "S",
		IsAdvanced:                 false,
		LaunchConfetti:             true,
		MaxNames:                   1000,
		PageBackgroundColor:        "#FFFFFF",
		PictureType:                "none",
		PlayClickWhenWinnerRemoved: false,
		ShowTitle:                  true,
		SlowSpin:                   false,
		SpinTime:                   10,
		Title:                      "Quests",
		Type:                       "color",
		WinnerMessage:              "",
	}

	return wh
}
