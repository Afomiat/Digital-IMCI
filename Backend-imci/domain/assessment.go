package domain

import (
	"time"
	"errors"
	"context"

	"github.com/google/uuid"
)

var (
	ErrAssessmentNotFound      = errors.New("assessment not found")
	ErrInvalidWeight           = errors.New("invalid weight")
	ErrInvalidAgeForAssessment = errors.New("patient age outside IMCI range (0-59 months)")
	ErrMedicalProfessionalAnswerNotFound = errors.New("medical professional answer not found")
)

type AssessmentType string
const (
	TypeYoungInfant AssessmentType = "young_infant"
	TypeChild       AssessmentType = "child"
)

type AssessmentStatus string
const (
	StatusDraft      AssessmentStatus = "draft"
	StatusInProgress AssessmentStatus = "in_progress"
	StatusClassified AssessmentStatus = "classified"
	StatusCompleted  AssessmentStatus = "completed"
	StatusCancelled  AssessmentStatus = "cancelled"
)

type JSONB map[string]interface{}

type Assessment struct {
	ID                    uuid.UUID         `json:"id"`
	MedicalProfessionalID uuid.UUID         `json:"medical_professional_id"`
	PatientID             uuid.UUID         `json:"patient_id"`
	AssessmentType        AssessmentType    `json:"assessment_type"`
	Status                AssessmentStatus  `json:"status"`
	WeightKg              float64           `json:"weight_kg"`
	Temperature           *float64          `json:"temperature,omitempty"`
	MainSymptoms          JSONB             `json:"main_symptoms"`
	MUAC                  *float64          `json:"muac,omitempty"`
	RespiratoryRate       *int              `json:"respiratory_rate,omitempty"`
	AgeMonths             int               `json:"age_months"`
	GuidelineVersion      string            `json:"guideline_version"`
	StartTime             time.Time         `json:"start_time"`
	EndTime               *time.Time        `json:"end_time,omitempty"`
	Summary               string            `json:"summary,omitempty"`
	IsOffline             bool              `json:"is_offline"`
	SyncedAt              *time.Time        `json:"synced_at,omitempty"`
	CreatedAt             time.Time         `json:"created_at"`
	UpdatedAt             time.Time         `json:"updated_at"`

	// NEW: Rule Engine Integration using existing columns
	OxygenSaturation      *int              `json:"oxygen_saturation,omitempty"`
	BilateralEdema        bool              `json:"bilateral_edema"`
	IsCriticalIllness     bool              `json:"is_critical_illness"`
	RequiresUrgentReferral bool             `json:"requires_urgent_referral"`
	HbLevel               *float64          `json:"hb_level,omitempty"`
    DevelopmentMilestones JSONB             `json:"development_milestones,omitempty"`
    JaundiceSigns         JSONB             `json:"jaundice_signs,omitempty"`

	MedicalProfessional *MedicalProfessional `json:"medical_professional,omitempty"`
	Patient             *Patient             `json:"patient,omitempty"`
	ClinicalFindings    *ClinicalFindings    `json:"clinical_findings,omitempty"`
	Classification      *Classification      `json:"classification,omitempty"`
	MedicalProfessionalAnswer *MedicalProfessionalAnswer `json:"medical_professional_answer,omitempty"`
}

// NEW: Medical Professional Answer represents our session
type MedicalProfessionalAnswer struct {
	ID                   uuid.UUID `json:"id"`
	AssessmentID         uuid.UUID `json:"assessment_id"`
	Answers              JSONB     `json:"answers"`
	QuestionSetVersion   string    `json:"question_set_version"`
	ClinicalFindings     JSONB     `json:"clinical_findings"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
// domain/clinical_findings.go
type ClinicalFindings struct {
	ID                    uuid.UUID `json:"id"`
	AssessmentID          uuid.UUID `json:"assessment_id"`
	
	// General Danger Signs (Existing)
	UnableToDrink         bool      `json:"unable_to_drink"`
	VomitsEverything      bool      `json:"vomits_everything"`
	HadConvulsions        bool      `json:"had_convulsions"`
	LethargicUnconscious  bool      `json:"lethargic_unconscious"`
	ConvulsingNow         bool      `json:"convulsing_now"`
	
	// Respiratory Assessment (Partially Existing)
	FastBreathing         bool      `json:"fast_breathing"`
	ChestIndrawing        bool      `json:"chest_indrawing"`
	Stridor               bool      `json:"stridor"`
	Wheezing              bool      `json:"wheezing"`
	OxygenSaturation      *int      `json:"oxygen_saturation"`
	RespiratoryRate       *int      `json:"respiratory_rate"` // NEW: Actual respiratory rate count
	CoughPresent          bool      `json:"cough_present"`    // NEW: Presence of cough
	CoughDurationDays     *int      `json:"cough_duration_days"` // NEW: How long cough has lasted
	
	// Diarrhea Assessment (Partially Existing)
	DiarrheaPresent       bool      `json:"diarrhea_present"` // NEW: Presence of diarrhea
	DiarrheaDurationDays  *int      `json:"diarrhea_duration_days"`
	BloodInStool          bool      `json:"blood_in_stool"`
	
	// Dehydration Signs (NEW)
	SunkenEyes            bool      `json:"sunken_eyes"`
	SkinPinchSlow         bool      `json:"skin_pinch_slow"`     // >2 seconds
	SkinPinchVerySlow     bool      `json:"skin_pinch_very_slow"` // Very slow skin pinch
	RestlessIrritable     bool      `json:"restless_irritable"`
	DrinkingEagerly       bool      `json:"drinking_eagerly"`    // For some dehydration
	DrinkingPoorly        bool      `json:"drinking_poorly"`     // For severe dehydration
	
	// Fever Assessment (Partially Existing)
	FeverPresent          bool      `json:"fever_present"` // NEW: Presence of fever
	FeverDurationDays     *int      `json:"fever_duration_days"`
	StiffNeck             bool      `json:"stiff_neck"`
	BulgingFontanelle     bool      `json:"bulging_fontanelle"`
	RunnyNose             bool      `json:"runny_nose"`           // NEW: For measles assessment
	RedEyes               bool      `json:"red_eyes"`             // NEW: For measles assessment
	GeneralizedRash       bool      `json:"generalized_rash"`     // NEW: For measles assessment
	
	// Measles (Existing)
	MeaslesNow            bool      `json:"measles_now"`
	MeaslesLast3Months    bool      `json:"measles_last_3_months"`
	
	// Ear Problems (NEW)
	EarPain               bool      `json:"ear_pain"`
	EarDischarge          bool      `json:"ear_discharge"`
	EarDischargeDurationDays *int   `json:"ear_discharge_duration_days"`
	TenderSwellingBehindEar bool    `json:"tender_swelling_behind_ear"`
	
	// Nutrition Assessment (Partially Existing)
	MUAC                  *float64  `json:"muac"`
	BilateralEdema        bool      `json:"bilateral_edema"`
	VisibleSevereWasting  bool      `json:"visible_severe_wasting"` // NEW: Visible signs
	WeightForHeightZScore *float64  `json:"weight_for_height_z_score"` // NEW: WFH z-score
	
	// Anemia Assessment (NEW)
	SeverePalmarPallor    bool      `json:"severe_palmar_pallor"`
	SomePalmarPallor      bool      `json:"some_palmar_pallor"`
	HbLevel               *float64  `json:"hb_level"`              // Hemoglobin level
	HctLevel              *float64  `json:"hct_level"`             // Hematocrit level
	
	// Jaundice Assessment (Existing)
	PalmsSolesYellow      bool      `json:"palms_soles_yellow"`
	SkinEyesYellow        bool      `json:"skin_eyes_yellow"`
	JaundiceAgeHours      *int      `json:"jaundice_age_hours"`
	
	// Young Infant Specific Signs (NEW)
	UnableToFeed          bool      `json:"unable_to_feed"`        // For infants <2 months
	NotFeedingWell        bool      `json:"not_feeding_well"`      // For infants <2 months
	MovementOnlyWhenStimulated bool `json:"movement_only_when_stimulated"`
	NoMovement            bool      `json:"no_movement"`
	UmbilicusRed          bool      `json:"umbilicus_red"`         // Local infection
	UmbilicusDrainingPus  bool      `json:"umbilicus_draining_pus"` // Local infection
	SkinPustules          bool      `json:"skin_pustules"`         // Local infection
	LowBodyTemperature    bool      `json:"low_body_temperature"`  // <35.5°C
	BodyTemperature       *float64  `json:"body_temperature"`      // Actual temperature
	
	// HIV/TB Assessment (NEW)
	HIVExposed            bool      `json:"hiv_exposed"`
	HIVStatusKnown        bool      `json:"hiv_status_known"`
	TBCoughDurationDays   *int      `json:"tb_cough_duration_days"` // >14 days
	TBContactHistory      bool      `json:"tb_contact_history"`
	TBWeightLoss          bool      `json:"tb_weight_loss"`
	NightSweats           bool      `json:"night_sweats"`
	
	// Development Assessment (NEW)
	SuspectedDevelopmentalDelay bool `json:"suspected_developmental_delay"`
	MilestonesAbsent     []string    `json:"milestones_absent"`     // List of absent milestones
	RiskFactorsPresent   []string    `json:"risk_factors_present"`  // Developmental risk factors
	
	// Feeding Assessment (NEW)
	Breastfeeding         bool      `json:"breastfeeding"`
	BreastfeedingFrequency *int     `json:"breastfeeding_frequency"` // Times in 24h
	ComplementaryFoods    bool      `json:"complementary_foods"`
	FeedingProblem        bool      `json:"feeding_problem"`
	Underweight           bool      `json:"underweight"`
	
	// Additional Clinical Notes
	OtherFindings         string    `json:"other_findings"`        // Free text for additional findings
	
	// Timestamps (Existing)
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// NEW: Classification (matches your table)
type Classification struct {
	ID                   uuid.UUID `json:"id"`
	AssessmentID         uuid.UUID `json:"assessment_id"`
	Disease              string    `json:"disease"`
	Color                string    `json:"color"`
	Details              string    `json:"details"`
	RuleVersion          string    `json:"rule_version"`
	ConfidenceScore      *float64  `json:"confidence_score,omitempty"`
	IsCriticalIllness    bool      `json:"is_critical_illness"`
	RequiresUrgentReferral bool    `json:"requires_urgent_referral"`
	TreatmentPriority    int       `json:"treatment_priority"`
	CreatedAt            time.Time `json:"created_at"`
}

// NEW: Treatment Plan (matches your table)
type TreatmentPlan struct {
	ID                  uuid.UUID `json:"id"`
	AssessmentID        uuid.UUID `json:"assessment_id"`
	ClassificationID    uuid.UUID `json:"classification_id"`
	DrugName            string    `json:"drug_name"`
	Dosage              string    `json:"dosage"`
	Frequency           string    `json:"frequency"`
	Duration            string    `json:"duration"`
	AdministrationRoute string    `json:"administration_route"`
	IsPreReferral       bool      `json:"is_pre_referral"`
	Instructions        string    `json:"instructions,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// Request/Response types
type CreateAssessmentRequest struct {
	PatientID       uuid.UUID `json:"patient_id" binding:"required"`
	WeightKg        float64   `json:"weight_kg" binding:"required,min=0.5,max=30.0"`
	Temperature     *float64  `json:"temperature,omitempty"`
	MainSymptoms    []string  `json:"main_symptoms,omitempty"`
	MUAC            *float64  `json:"muac,omitempty"`
	RespiratoryRate *int      `json:"respiratory_rate,omitempty"`
	IsOffline       bool      `json:"is_offline"`
}

type StartAssessmentFlowRequest struct {
	AssessmentID uuid.UUID `json:"assessment_id" binding:"required"`
}

type SubmitAnswerRequest struct {
	AssessmentID uuid.UUID   `json:"assessment_id" binding:"required"`
	NodeID       string      `json:"node_id" binding:"required"`
	Answer       interface{} `json:"answer" binding:"required"`
}

type AssessmentFlowResponse struct {
	SessionID       uuid.UUID              `json:"session_id"`
	Question        *FlowQuestion          `json:"question,omitempty"`
	Classification  *Classification        `json:"classification,omitempty"`
	ClinicalFindings *ClinicalFindings     `json:"clinical_findings,omitempty"`
	IsComplete      bool                   `json:"is_complete"`
}

type FlowQuestion struct {
	NodeID      string    `json:"node_id"`
	Type        string    `json:"type"`
	Question    string    `json:"question"`
	Instruction string    `json:"instruction,omitempty"`
	Options     []Option  `json:"options,omitempty"`
}

type Option struct {
	Value string `json:"value"`
	Text  string `json:"text"`
}
type UpdateAssessmentRequest struct {
    // Example fields (replace with actual updatable fields)
    Diagnosis   string  `json:"diagnosis,omitempty"`
    Notes       string  `json:"notes,omitempty"`
    Weight      float64 `json:"weight,omitempty"`
    Temperature float64 `json:"temperature,omitempty"`
    // Add other fields as necessary
}
// NEW: Repository interfaces for your tables
type MedicalProfessionalAnswerRepository interface {
	Create(ctx context.Context, answer *MedicalProfessionalAnswer) error
	GetByAssessmentID(ctx context.Context, assessmentID uuid.UUID) (*MedicalProfessionalAnswer, error)
	Update(ctx context.Context, answer *MedicalProfessionalAnswer) error
	Upsert(ctx context.Context, answer *MedicalProfessionalAnswer) error
}

type ClinicalFindingsRepository interface {
	Create(ctx context.Context, findings *ClinicalFindings) error
	GetByAssessmentID(ctx context.Context, assessmentID uuid.UUID) (*ClinicalFindings, error)
	Upsert(ctx context.Context, findings *ClinicalFindings) error
}

type ClassificationRepository interface {
	Create(ctx context.Context, classification *Classification) error
	GetByAssessmentID(ctx context.Context, assessmentID uuid.UUID) (*Classification, error)
	Upsert(ctx context.Context, classification *Classification) error
}

type TreatmentPlanRepository interface {
	Create(ctx context.Context, plan *TreatmentPlan) error
	GetByAssessmentID(ctx context.Context, assessmentID uuid.UUID) ([]*TreatmentPlan, error)
}

type AssessmentRepository interface {
	Create(ctx context.Context, assessment *Assessment) error
	GetByID(ctx context.Context, id uuid.UUID, medicalProfessionalID uuid.UUID) (*Assessment, error)
	GetByPatientID(ctx context.Context, patientID uuid.UUID, medicalProfessionalID uuid.UUID) ([]*Assessment, error)
	Update(ctx context.Context, assessment *Assessment) error
	Delete(ctx context.Context, id uuid.UUID, medicalProfessionalID uuid.UUID) error
	CalculateAgeInfo(ctx context.Context, patientID uuid.UUID, assessmentTime time.Time) (int, AssessmentType, error)
}

type AssessmentUsecase interface {
	CreateAssessment(ctx context.Context, req *CreateAssessmentRequest, medicalProfessionalID uuid.UUID) (*Assessment, error)
	GetAssessment(ctx context.Context, assessmentID uuid.UUID, medicalProfessionalID uuid.UUID) (*Assessment, error)
	GetAssessmentsByPatient(ctx context.Context, patientID uuid.UUID, medicalProfessionalID uuid.UUID) ([]*Assessment, error)
	UpdateAssessment(ctx context.Context, assessment *Assessment) error
	DeleteAssessment(ctx context.Context, assessmentID uuid.UUID, medicalProfessionalID uuid.UUID) error
}