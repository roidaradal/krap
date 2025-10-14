# Krap
Kollection of Roi's Application Patterns

## audit 
Audit Logs

* audit.Initialize( ) error
* audit.NewDetails(...string) string 
* audit.NewActionLog(actorID, action, details string) *ActionLog
* audit.NewActionLogs(actorID, [][2]string) []*ActionLog 
* audit.AddActionLogTx(*ze.Request, *ActionLog, table string) error 
* audit.AddActionLogsTx(*ze.Request, []*ActionLog, table string) error
* audit.NewBatchLog(action, details, actionGlue string) *BatchLog 
* audit.NewBatchLogItems(batchCode string, detailsList []string) []*BatchLogItem
* audit.AddBatchLogTx(*ze.Request, *BatchLog) error 
* audit.AddBatchLogItemsTx(*ze.Request, []*BatchLogItem) error

## authn
Authentication

## authz
Authorization

* authz.Initialize( ) error 
* authz.LoadAccess(*ze.Request) error 
* authz.LoadScopedAccess(*ze.Request, table string) error 
* authz.GetAllAccess( ) dict.StringListMap 
* authz.GetAllRoleAccess( ) dict.StringListMap 
* authz.GetScopedAccess(table, scopeCode string) dict.StringListMap
* authz.CheckActionAllowedFor(*ze.Request, action, item, role string) error 
* authz.CheckScopedActionAllowedFor(*ze.Request, table, scopeCode, action, item, role string) error

## config 
Configuration and Features

* config.Initialize( ) error
* config.Lookup(*ze.Request, appKeys []string) (dict.StringMap, error)
* config.Create[T any](*T, dict.StringMap, *Defaults) *T
* config.LoadFeatures(*ze.Request) error 
* config.LoadScopedFeatures(*ze.Request, table string) error 
* config.GetAllFeatures( ) dict.BoolMap 
* config.GetActiveFeatures( ) []string 
* config.GetAllScopedFeatures(table string) dict.StringListMap 
* config.GetAllFeatureScopes(table string) dict.StringListMap
* config.GetScopedFeatures(table string, scopeCodes ...string) dict.StringListMap
* config.CheckFeature(*ze.Request, feature string) error 
* config.CheckScopedFeature(*ze.Request, table, scope, feature string) error

## daemon
Daemons

* daemon.LoadConfig[T any](path string) (*T, error)
* daemon.Run(name string, task func(), interval int, time.Duration)