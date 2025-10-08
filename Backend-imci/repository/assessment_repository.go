// repository/assessment_repo.go
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AssessmentRepo struct {
	db *pgxpool.Pool
}

func NewAssessmentRepo(db *pgxpool.Pool) domain.AssessmentRepository {
	return &AssessmentRepo{db: db}
}

func (r *AssessmentRepo) Create(ctx context.Context, assessment *domain.Assessment) error {
	query := `
		INSERT INTO assessments (
			id, medical_professional_id, patient_id, assessment_type, status, 
			weight_kg, temperature, main_symptoms, muac, respiratory_rate,
			age_months, guideline_version, start_time, is_offline, created_at, updated_at,
			oxygen_saturation, jaundice_signs, development_milestones, hb_level,
			bilateral_edema, is_critical_illness, requires_urgent_referral
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
	`

	mainSymptomsJSON, err := json.Marshal(assessment.MainSymptoms)
	if err != nil {
		return fmt.Errorf("failed to marshal main symptoms: %w", err)
	}

	jaundiceSignsJSON, err := json.Marshal(assessment.JaundiceSigns)
	if err != nil {
		return fmt.Errorf("failed to marshal jaundice signs: %w", err)
	}

	developmentMilestonesJSON, err := json.Marshal(assessment.DevelopmentMilestones)
	if err != nil {
		return fmt.Errorf("failed to marshal development milestones: %w", err)
	}

	now := time.Now()
	assessment.CreatedAt = now
	assessment.UpdatedAt = now

	_, err = r.db.Exec(ctx, query,
		assessment.ID,
		assessment.MedicalProfessionalID,
		assessment.PatientID,
		assessment.AssessmentType,
		assessment.Status,
		assessment.WeightKg,
		assessment.Temperature,
		mainSymptomsJSON,
		assessment.MUAC,
		assessment.RespiratoryRate,
		assessment.AgeMonths,
		assessment.GuidelineVersion,
		assessment.StartTime,
		assessment.IsOffline,
		assessment.CreatedAt,
		assessment.UpdatedAt,
		assessment.OxygenSaturation,
		jaundiceSignsJSON,
		developmentMilestonesJSON,
		assessment.HbLevel,
		assessment.BilateralEdema,
		assessment.IsCriticalIllness,
		assessment.RequiresUrgentReferral,
	)

	if err != nil {
		return fmt.Errorf("failed to create assessment: %w", err)
	}

	return nil
}

func (r *AssessmentRepo) GetByID(ctx context.Context, id uuid.UUID, medicalProfessionalID uuid.UUID) (*domain.Assessment, error) {
	query := `
		SELECT 
			id, medical_professional_id, patient_id, assessment_type, status,
			weight_kg, temperature, main_symptoms, muac, respiratory_rate,
			age_months, guideline_version, start_time, end_time, summary,
			is_offline, synced_at, created_at, updated_at,
			oxygen_saturation, jaundice_signs, development_milestones, hb_level,
			bilateral_edema, is_critical_illness, requires_urgent_referral
		FROM assessments 
		WHERE id = $1 AND medical_professional_id = $2
	`

	var assessment domain.Assessment
	var mainSymptoms, jaundiceSigns, developmentMilestones []byte
	var temperature, muac, hbLevel sql.NullFloat64
	var respiratoryRate, oxygenSaturation sql.NullInt32
	var endTime, syncedAt sql.NullTime
	var summary sql.NullString

	// Add this before the Scan to see the exact query
fmt.Printf("🔍 Executing query with params: assessmentID=%s, medicalProfessionalID=%s\n", id, medicalProfessionalID)

// Count the number of ? placeholders in your SELECT vs Scan arguments
selectCount := strings.Count(query, ",") + 1 // +1 for the first column
fmt.Printf("🔍 SELECT has %d columns, Scan has 26 destinations\n", selectCount)

	err := r.db.QueryRow(ctx, query, id, medicalProfessionalID).Scan(
		&assessment.ID,                       // 1
		&assessment.MedicalProfessionalID,    // 2
		&assessment.PatientID,                // 3
		&assessment.AssessmentType,           // 4
		&assessment.Status,                   // 5
		&assessment.WeightKg,                 // 6
		&temperature,                         // 7
		&mainSymptoms,                        // 8
		&muac,                                // 9
		&respiratoryRate,                     // 10
		&assessment.AgeMonths,                // 11
		&assessment.GuidelineVersion,         // 12
		&assessment.StartTime,                // 13
		&endTime,                             // 14
		&summary,                             // 15
		&assessment.IsOffline,                // 16
		&syncedAt,                            // 17
		&assessment.CreatedAt,                // 18
		&assessment.UpdatedAt,                // 19
		&oxygenSaturation,                    // 20
		&jaundiceSigns,                       // 21
		&developmentMilestones,               // 22
		&hbLevel,                             // 23
		&assessment.BilateralEdema,           // 24
		&assessment.IsCriticalIllness,        // 25
		&assessment.RequiresUrgentReferral,   // 26
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrAssessmentNotFound
		}
		return nil, fmt.Errorf("failed to get assessment: %w", err)
	}

	// Handle nullable fields
	if temperature.Valid {
		temp := temperature.Float64
		assessment.Temperature = &temp
	}
	if muac.Valid {
		muacVal := muac.Float64
		assessment.MUAC = &muacVal
	}
	if respiratoryRate.Valid {
		rate := int(respiratoryRate.Int32)
		assessment.RespiratoryRate = &rate
	}
	if oxygenSaturation.Valid {
		sat := int(oxygenSaturation.Int32)
		assessment.OxygenSaturation = &sat
	}
	if hbLevel.Valid {
		hb := hbLevel.Float64
		assessment.HbLevel = &hb
	}
	if endTime.Valid {
		assessment.EndTime = &endTime.Time
	}
	if syncedAt.Valid {
		assessment.SyncedAt = &syncedAt.Time
	}
	if summary.Valid {
		assessment.Summary = summary.String
	}

	// Unmarshal JSONB fields
	if err := json.Unmarshal(mainSymptoms, &assessment.MainSymptoms); err != nil {
		return nil, fmt.Errorf("failed to unmarshal main symptoms: %w", err)
	}
	if err := json.Unmarshal(jaundiceSigns, &assessment.JaundiceSigns); err != nil {
		return nil, fmt.Errorf("failed to unmarshal jaundice signs: %w", err)
	}
	if err := json.Unmarshal(developmentMilestones, &assessment.DevelopmentMilestones); err != nil {
		return nil, fmt.Errorf("failed to unmarshal development milestones: %w", err)
	}

	return &assessment, nil
}

func (r *AssessmentRepo) Update(ctx context.Context, assessment *domain.Assessment) error {
	query := `
		UPDATE assessments 
		SET status = $1, temperature = $2, main_symptoms = $3, muac = $4, 
			respiratory_rate = $5, end_time = $6, summary = $7, synced_at = $8,
			updated_at = $9, oxygen_saturation = $10, jaundice_signs = $11,
			development_milestones = $12, hb_level = $13, bilateral_edema = $14,
			is_critical_illness = $15, requires_urgent_referral = $16
		WHERE id = $17 AND medical_professional_id = $18
	`

	mainSymptomsJSON, err := json.Marshal(assessment.MainSymptoms)
	if err != nil {
		return fmt.Errorf("failed to marshal main symptoms: %w", err)
	}

	jaundiceSignsJSON, err := json.Marshal(assessment.JaundiceSigns)
	if err != nil {
		return fmt.Errorf("failed to marshal jaundice signs: %w", err)
	}

	developmentMilestonesJSON, err := json.Marshal(assessment.DevelopmentMilestones)
	if err != nil {
		return fmt.Errorf("failed to marshal development milestones: %w", err)
	}

	assessment.UpdatedAt = time.Now()

	result, err := r.db.Exec(ctx, query,
		assessment.Status,
		assessment.Temperature,
		mainSymptomsJSON,
		assessment.MUAC,
		assessment.RespiratoryRate,
		assessment.EndTime,
		assessment.Summary,
		assessment.SyncedAt,
		assessment.UpdatedAt,
		assessment.OxygenSaturation,
		jaundiceSignsJSON,
		developmentMilestonesJSON,
		assessment.HbLevel,
		assessment.BilateralEdema,
		assessment.IsCriticalIllness,
		assessment.RequiresUrgentReferral,
		assessment.ID,
		assessment.MedicalProfessionalID,
	)

	if err != nil {
		return fmt.Errorf("failed to update assessment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrAssessmentNotFound
	}

	return nil
}

func (r *AssessmentRepo) CalculateAgeInfo(ctx context.Context, patientID uuid.UUID, assessmentTime time.Time) (int, domain.AssessmentType, error) {
	query := `
		SELECT 
			EXTRACT(YEAR FROM AGE($1, date_of_birth)) * 12 + 
			EXTRACT(MONTH FROM AGE($1, date_of_birth)) as age_months
		FROM patients 
		WHERE id = $2
	`

	var ageMonths int
	err := r.db.QueryRow(ctx, query, assessmentTime, patientID).Scan(&ageMonths)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, "", domain.ErrPatientNotFound
		}
		return 0, "", fmt.Errorf("failed to calculate age: %w", err)
	}

	// Validate age range (IMCI: 0-59 months)
	if ageMonths < 0 || ageMonths > 59 {
		return 0, "", domain.ErrInvalidAgeForAssessment
	}

	// Determine assessment type (IMCI: <2 months = young infant)
	var assessmentType domain.AssessmentType
	if ageMonths < 2 {
		assessmentType = domain.TypeYoungInfant
	} else {
		assessmentType = domain.TypeChild
	}

	return ageMonths, assessmentType, nil
}

func (r *AssessmentRepo) GetByPatientID(ctx context.Context, patientID uuid.UUID, medicalProfessionalID uuid.UUID) ([]*domain.Assessment, error) {
	query := `
		SELECT id, medical_professional_id, patient_id, assessment_type, status,
			weight_kg, temperature, main_symptoms, muac, respiratory_rate,
			age_months, guideline_version, start_time, end_time, summary,
			is_offline, synced_at, created_at, updated_at
		FROM assessments 
		WHERE patient_id = $1 AND medical_professional_id = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, patientID, medicalProfessionalID)
	if err != nil {
		return nil, fmt.Errorf("failed to query assessments: %w", err)
	}
	defer rows.Close()

	var assessments []*domain.Assessment
	for rows.Next() {
		var assessment domain.Assessment
		var mainSymptoms []byte
		var temperature, muac sql.NullFloat64
		var respiratoryRate sql.NullInt32
		var endTime, syncedAt sql.NullTime
		var summary sql.NullString

		err := rows.Scan(
			&assessment.ID,
			&assessment.MedicalProfessionalID,
			&assessment.PatientID,
			&assessment.AssessmentType,
			&assessment.Status,
			&assessment.WeightKg,
			&temperature,
			&mainSymptoms,
			&muac,
			&respiratoryRate,
			&assessment.AgeMonths,
			&assessment.GuidelineVersion,
			&assessment.StartTime,
			&endTime,
			&summary,
			&assessment.Summary,
			&assessment.IsOffline,
			&syncedAt,
			&assessment.CreatedAt,
			&assessment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan assessment: %w", err)
		}

		// Handle nullable fields
		if temperature.Valid {
			temp := temperature.Float64
			assessment.Temperature = &temp
		}
		if muac.Valid {
			muacVal := muac.Float64
			assessment.MUAC = &muacVal
		}
		if respiratoryRate.Valid {
			rate := int(respiratoryRate.Int32)
			assessment.RespiratoryRate = &rate
		}
		if endTime.Valid {
			assessment.EndTime = &endTime.Time
		}
		if syncedAt.Valid {
			assessment.SyncedAt = &syncedAt.Time
		}
		if summary.Valid { // ✅ Handle NULL summary
			assessment.Summary = summary.String
		}

		// Unmarshal JSONB
		if err := json.Unmarshal(mainSymptoms, &assessment.MainSymptoms); err != nil {
			return nil, fmt.Errorf("failed to unmarshal main symptoms: %w", err)
		}

		assessments = append(assessments, &assessment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating assessments: %w", err)
	}

	return assessments, nil
}

func (r *AssessmentRepo) Delete(ctx context.Context, id uuid.UUID, medicalProfessionalID uuid.UUID) error {
	query := `DELETE FROM assessments WHERE id = $1 AND medical_professional_id = $2`

	result, err := r.db.Exec(ctx, query, id, medicalProfessionalID)
	if err != nil {
		return fmt.Errorf("failed to delete assessment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrAssessmentNotFound
	}

	return nil
}