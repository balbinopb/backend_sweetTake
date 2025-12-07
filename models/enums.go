package models

type RoleType string
const (
    RoleUser         RoleType = "User"
    RoleAdministrator RoleType = "Administrator"
)

type AuthLevelType string
const (
    AuthBasic AuthLevelType = "basic"
    AuthAdmin AuthLevelType = "admin"
)

type PreferenceType string
const (
    PrefLowSugar      PreferenceType = "Low Sugar Diet"
    PrefDiabetic      PreferenceType = "Diabetic-Friendly Diet"
    PrefBalanced      PreferenceType = "Balanced Diet"
)

type HealthGoalType string
const (
    GoalReduceSugar   HealthGoalType = "Reduce Daily Sugar Intake"
    GoalMaintainLevel HealthGoalType = "Maintain Stable Blood Sugar"
    GoalWeightLoss    HealthGoalType = "Weight Loss"
)

type ReadStatusType string
const (
    ReadStatusRead   ReadStatusType = "read"
    ReadStatusUnread ReadStatusType = "unread"
)

type RiskLevelType string
const (
    RiskLow      RiskLevelType = "Low"
    RiskModerate RiskLevelType = "Moderate"
    RiskHigh     RiskLevelType = "High"
)
