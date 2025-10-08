package domain

import (
	"time"

	"github.com/google/uuid"
)

type NodeType string

const (
	NodeTypeMultipleChoice NodeType = "multiple_choice"
	NodeTypeYesNo          NodeType = "yes_no" 
	NodeTypeClassification NodeType = "classification"
	NodeTypeAssessment     NodeType = "assessment"
)

type DecisionNode struct {
	ID          string            `json:"id"`
	Type        NodeType          `json:"type"`
	Question    string            `json:"question,omitempty"`
	Instruction string            `json:"instruction,omitempty"`
	Options     []Option          `json:"options,omitempty"`
	Rules       []Rule            `json:"rules,omitempty"`
	Classification *Classification `json:"classification,omitempty"`
	Assessment  string            `json:"assessment,omitempty"` // e.g., "respiratory_rate"
}

type Option struct {
	Value string `json:"value"`
	Text  string `json:"text"`
}

type Rule struct {
	Condition       string   `json:"condition"`
	SelectedOptions []string `json:"selected_options,omitempty"`
	NextNode        string   `json:"next_node"`
	Thresholds      map[string]int `json:"thresholds,omitempty"` // For respiratory rate
	Value           string            `json:"value,omitempty"`  // ← ADD THIS FIELD

}

type Classification struct {
	Color      string        `json:"color"`
	Name       string        `json:"name"`
	TreatmentPlan *TreatmentPlan `json:"treatment_plan"`
}

type TreatmentPlan struct {
	RequiresReferral      bool     `json:"requires_referral"`
	Urgency              string   `json:"urgency,omitempty"`
	PreReferralTreatments []string `json:"pre_referral_treatments,omitempty"`
	Drugs                []Drug   `json:"drugs,omitempty"`
}

type Drug struct {
	Name     string `json:"name"`
	Dosage   string `json:"dosage"`
	Duration string `json:"duration"`
}

type DecisionTree struct {
	StartNode string                 `json:"start_node"`
	Nodes     map[string]DecisionNode `json:"nodes"`
}

type AssessmentSession struct {
    SessionID        uuid.UUID                `json:"session_id"`
    AssessmentID     uuid.UUID                `json:"assessment_id"`
    AssessmentType   string                   `json:"assessment_type"` // "child" or "young_infant"
    CurrentNodeID    string                   `json:"current_node_id"`
    Answers          map[string]interface{}   `json:"answers"`
    ClinicalFindings *ClinicalFindings        `json:"clinical_findings"`
    Classification   *Classification          `json:"classification"`
    CreatedAt        time.Time                `json:"created_at"`
    UpdatedAt        time.Time                `json:"updated_at"`
}

// Use the same ClinicalFindings struct as your main domain
type ClinicalFindings struct {
	ID                    uuid.UUID `json:"id"`
	AssessmentID          uuid.UUID `json:"assessment_id"`
	UnableToDrink         bool      `json:"unable_to_drink"`
	VomitsEverything      bool      `json:"vomits_everything"`
	HadConvulsions        bool      `json:"had_convulsions"`
	LethargicUnconscious  bool      `json:"lethargic_unconscious"`
	ConvulsingNow         bool      `json:"convulsing_now"`
	FastBreathing         bool      `json:"fast_breathing"`
	ChestIndrawing        bool      `json:"chest_indrawing"`
	Stridor               bool      `json:"stridor"`
	Wheezing              bool      `json:"wheezing"`
	OxygenSaturation      *int      `json:"oxygen_saturation"`
	RespiratoryRate       *int      `json:"respiratory_rate"`
	CoughPresent          bool      `json:"cough_present"`
	CoughDurationDays     *int      `json:"cough_duration_days"`
	DiarrheaPresent       bool      `json:"diarrhea_present"`
	DiarrheaDurationDays  *int      `json:"diarrhea_duration_days"`
	BloodInStool          bool      `json:"blood_in_stool"`
	SunkenEyes            bool      `json:"sunken_eyes"`
	SkinPinchSlow         bool      `json:"skin_pinch_slow"`
	SkinPinchVerySlow     bool      `json:"skin_pinch_very_slow"`
	RestlessIrritable     bool      `json:"restless_irritable"`
	DrinkingEagerly       bool      `json:"drinking_eagerly"`
	DrinkingPoorly        bool      `json:"drinking_poorly"`
	FeverPresent          bool      `json:"fever_present"`
	FeverDurationDays     *int      `json:"fever_duration_days"`
	StiffNeck             bool      `json:"stiff_neck"`
	BulgingFontanelle     bool      `json:"bulging_fontanelle"`
	RunnyNose             bool      `json:"runny_nose"`
	RedEyes               bool      `json:"red_eyes"`
	GeneralizedRash       bool      `json:"generalized_rash"`
	MeaslesNow            bool      `json:"measles_now"`
	MeaslesLast3Months    bool      `json:"measles_last_3_months"`
	EarPain               bool      `json:"ear_pain"`
	EarDischarge          bool      `json:"ear_discharge"`
	EarDischargeDurationDays *int   `json:"ear_discharge_duration_days"`
	TenderSwellingBehindEar bool    `json:"tender_swelling_behind_ear"`
	MUAC                  *float64  `json:"muac"`
	BilateralEdema        bool      `json:"bilateral_edema"`
	VisibleSevereWasting  bool      `json:"visible_severe_wasting"`
	WeightForHeightZScore *float64  `json:"weight_for_height_z_score"`
	SeverePalmarPallor    bool      `json:"severe_palmar_pallor"`
	SomePalmarPallor      bool      `json:"some_palmar_pallor"`
	HbLevel               *float64  `json:"hb_level"`
	HctLevel              *float64  `json:"hct_level"`
	PalmsSolesYellow      bool      `json:"palms_soles_yellow"`
	SkinEyesYellow        bool      `json:"skin_eyes_yellow"`
	JaundiceAgeHours      *int      `json:"jaundice_age_hours"`
	UnableToFeed          bool      `json:"unable_to_feed"`
	NotFeedingWell        bool      `json:"not_feeding_well"`
	MovementOnlyWhenStimulated bool `json:"movement_only_when_stimulated"`
	NoMovement            bool      `json:"no_movement"`
	UmbilicusRed          bool      `json:"umbilicus_red"`
	UmbilicusDrainingPus  bool      `json:"umbilicus_draining_pus"`
	SkinPustules          bool      `json:"skin_pustules"`
	LowBodyTemperature    bool      `json:"low_body_temperature"`
	BodyTemperature       *float64  `json:"body_temperature"`
	HIVExposed            bool      `json:"hiv_exposed"`
	HIVStatusKnown        bool      `json:"hiv_status_known"`
	TBCoughDurationDays   *int      `json:"tb_cough_duration_days"`
	TBContactHistory      bool      `json:"tb_contact_history"`
	TBWeightLoss          bool      `json:"tb_weight_loss"`
	NightSweats           bool      `json:"night_sweats"`
	SuspectedDevelopmentalDelay bool `json:"suspected_developmental_delay"`
	MilestonesAbsent     []string    `json:"milestones_absent"`
	RiskFactorsPresent   []string    `json:"risk_factors_present"`
	Breastfeeding         bool      `json:"breastfeeding"`
	BreastfeedingFrequency *int     `json:"breastfeeding_frequency"`
	ComplementaryFoods    bool      `json:"complementary_foods"`
	FeedingProblem        bool      `json:"feeding_problem"`
	Underweight           bool      `json:"underweight"`
	OtherFindings         string    `json:"other_findings"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}