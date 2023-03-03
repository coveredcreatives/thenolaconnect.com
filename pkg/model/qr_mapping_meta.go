package model

import "gorm.io/gorm"

type QRMappingMeta struct {
	QREncodedData     string `json:"qr_encoded_data"`
	UniqueImpressions int    `json:"unique_impressions"`
}

func CalculateUniqueImpressions(db *gorm.DB, pkeys []string) map[string]QRMappingMeta {
	qr_mapping_metas_by_id := map[string]QRMappingMeta{}
	qr_mapping_metas := []QRMappingMeta{}
	if len(pkeys) == 0 {
		return qr_mapping_metas_by_id
	}

	db.
		Raw(`SELECT
		qr_mapping_impression.qr_encoded_data as qr_encoded_data,
		count(distinct qr_mapping_impression.qr_encoded_data) as unique_impressions
		FROM qr_mapping.qr_mapping_impression
		WHERE qr_mapping_impression.qr_encoded_data IN ?
		GROUP BY qr_mapping_impression.qr_encoded_data`, pkeys).
		Scan(&qr_mapping_metas)
	for _, meta := range qr_mapping_metas {
		qr_mapping_metas_by_id[meta.QREncodedData] = meta
	}
	return qr_mapping_metas_by_id
}
