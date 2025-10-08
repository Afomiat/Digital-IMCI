// repository/clinical_findings_repo.go
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClinicalFindingsRepo struct {
	db *pgxpool.Pool
}

func NewClinicalFindingsRepo(db *pgxpool.Pool) domain.ClinicalFindingsRepository {
	return &ClinicalFindingsRepo{db: db}
}

func (r *ClinicalFindingsRepo) Create(ctx context.Context, findings *domain.ClinicalFindings) error {
	query := `
		INSERT INTO clinical_findings (
			id, assessment_id, unable_to_drink, vomits_everything, had_convulsions,
			lethargic_unconscious, convulsing_now, fast_breathing, chest_indrawing,
			stridor, wheezing, oxygen_saturation, palms_soles_yellow, skin_eyes_yellow,
			jaundice_age_hours, diarrhea_duration_days, blood_in_stool, fever_duration_days,
			measles_now, measles_last_3_months, stiff_neck, bulging_fontanelle, muac,
			bilateral_edema, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26)
	`

	now := time.Now()
	findings.CreatedAt = now
	findings.UpdatedAt = now

	_, err := r.db.Exec(ctx, query,
		findings.ID,
		findings.AssessmentID,
		findings.UnableToDrink,
		findings.VomitsEverything,
		findings.HadConvulsions,
		findings.LethargicUnconscious,
		findings.ConvulsingNow,
		findings.FastBreathing,
		findings.ChestIndrawing,
		findings.Stridor,
		findings.Wheezing,
		findings.OxygenSaturation,
		findings.PalmsSolesYellow,
		findings.SkinEyesYellow,
		findings.JaundiceAgeHours,
		findings.DiarrheaDurationDays,
		findings.BloodInStool,
		findings.FeverDurationDays,
		findings.MeaslesNow,
		findings.MeaslesLast3Months,
		findings.StiffNeck,
		findings.BulgingFontanelle,
		findings.MUAC,
		findings.BilateralEdema,
		findings.CreatedAt,
		findings.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create clinical findings: %w", err)
	}

	return nil
}

func (r *ClinicalFindingsRepo) GetByAssessmentID(ctx context.Context, assessmentID uuid.UUID) (*domain.ClinicalFindings, error) {
	query := `
		SELECT id, assessment_id, unable_to_drink, vomits_everything, had_convulsions,
			lethargic_unconscious, convulsing_now, fast_breathing, chest_indrawing,
			stridor, wheezing, oxygen_saturation, palms_soles_yellow, skin_eyes_yellow,
			jaundice_age_hours, diarrhea_duration_days, blood_in_stool, fever_duration_days,
			measles_now, measles_last_3_months, stiff_neck, bulging_fontanelle, muac,
			bilateral_edema, created_at, updated_at
		FROM clinical_findings 
		WHERE assessment_id = $1
	`

	var findings domain.ClinicalFindings
	err := r.db.QueryRow(ctx, query, assessmentID).Scan(
		&findings.ID,
		&findings.AssessmentID,
		&findings.UnableToDrink,
		&findings.VomitsEverything,
		&findings.HadConvulsions,
		&findings.LethargicUnconscious,
		&findings.ConvulsingNow,
		&findings.FastBreathing,
		&findings.ChestIndrawing,
		&findings.Stridor,
		&findings.Wheezing,
		&findings.OxygenSaturation,
		&findings.PalmsSolesYellow,
		&findings.SkinEyesYellow,
		&findings.JaundiceAgeHours,
		&findings.DiarrheaDurationDays,
		&findings.BloodInStool,
		&findings.FeverDurationDays,
		&findings.MeaslesNow,
		&findings.MeaslesLast3Months,
		&findings.StiffNeck,
		&findings.BulgingFontanelle,
		&findings.MUAC,
		&findings.BilateralEdema,
		&findings.CreatedAt,
		&findings.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("clinical findings not found for assessment: %w", err)
		}
		return nil, fmt.Errorf("failed to get clinical findings: %w", err)
	}

	return &findings, nil
}

func (r *ClinicalFindingsRepo) Upsert(ctx context.Context, findings *domain.ClinicalFindings) error {
	existing, err := r.GetByAssessmentID(ctx, findings.AssessmentID)
	if err != nil {
		// Create new if not exists
		return r.Create(ctx, findings)
	}
	
	// Update existing
	findings.ID = existing.ID
	return r.Update(ctx, findings)
}

func (r *ClinicalFindingsRepo) Update(ctx context.Context, findings *domain.ClinicalFindings) error {
	query := `
		UPDATE clinical_findings 
		SET unable_to_drink = $1, vomits_everything = $2, had_convulsions = $3,
			lethargic_unconscious = $4, convulsing_now = $5, fast_breathing = $6,
			chest_indrawing = $7, stridor = $8, wheezing = $9, oxygen_saturation = $10,
			palms_soles_yellow = $11, skin_eyes_yellow = $12, jaundice_age_hours = $13,
			diarrhea_duration_days = $14, blood_in_stool = $15, fever_duration_days = $16,
			measles_now = $17, measles_last_3_months = $18, stiff_neck = $19,
			bulging_fontanelle = $20, muac = $21, bilateral_edema = $22, updated_at = $23
		WHERE id = $24
	`

	findings.UpdatedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		findings.UnableToDrink,
		findings.VomitsEverything,
		findings.HadConvulsions,
		findings.LethargicUnconscious,
		findings.ConvulsingNow,
		findings.FastBreathing,
		findings.ChestIndrawing,
		findings.Stridor,
		findings.Wheezing,
		findings.OxygenSaturation,
		findings.PalmsSolesYellow,
		findings.SkinEyesYellow,
		findings.JaundiceAgeHours,
		findings.DiarrheaDurationDays,
		findings.BloodInStool,
		findings.FeverDurationDays,
		findings.MeaslesNow,
		findings.MeaslesLast3Months,
		findings.StiffNeck,
		findings.BulgingFontanelle,
		findings.MUAC,
		findings.BilateralEdema,
		findings.UpdatedAt,
		findings.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update clinical findings: %w", err)
	}

	return nil
}