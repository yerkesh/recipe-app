package constant

import "fmt"

const (
	MsgCreated  = "Объект создан"
	MsgUpdated  = "Успешно обновлен"
	MsgDeleted  = "Успешно удален"
	MsgPatched  = "Успешно частично обновлен"
	MsgSuccess  = "Успешно"
	MsgRequired = "required"
)

// Request errors.
const (
	MsgNotFoundErr    = "Не найдена запись в БД"
	MsgRequiredErr    = "Не отправлены обязательные поля"
	MsgUnhandledErr   = "Непредвиденная ошибка"
	MsgRequestBodyErr = "Переданы некорректные данные"
	MsgAuthorizeErr   = "Ошибка авторизации"
	MsgAlreadyExists  = "Такая запись уже существует в БД"
)

// Tablenames.
type Table string

const (
	TblPopularity Table = "popularity"
	TblCommodity  Table = "commodity"
	TblEquipment  Table = "equipment"

	TblWarehouse               Table = "warehouse"
	TblJuncWarehousePopularity Table = "warehouse_popularity"
	TblJuncWarehouseCommodity  Table = "warehouse_commodity"
	TblJuncWarehouseEquipment  Table = "warehouse_equipment"
	TblAddress                 Table = "address"
	TblFunctionalZone          Table = "functional_zone"

	TblErpWarehouse Table = "erp_warehouse"

	TblSystemParam               Table = "system_param"
	TblFunctionalZoneSystemParam Table = "functional_zone_system_param"
	TblJuncWarehouseSystemParam  Table = "warehouse_system_param"

	TblCellType           Table = "cell_type"
	TblCell               Table = "cell"
	TblCellMeasurement    Table = "cell_measurement"
	TblCellGeneral        Table = "cell_general"
	TblJuncCellPopularity Table = "cell_popularity"
	TblJuncCellCommodity  Table = "cell_commodity"
	TblJuncCellEquipment  Table = "cell_equipment"

	TblRamp     Table = "ramp"
	TblRampLine Table = "ramp_line"

	TblCellItem Table = "cell_item"
)

type CEPTable Table

const (
	TblCEPPopularity = CEPTable(TblPopularity)
	TblCEPCommodity  = CEPTable(TblCommodity)
	TblCEPEquipment  = CEPTable(TblEquipment)
)

func (t Table) As(as ...string) string {
	if len(as) == 0 {
		return t.String()
	}

	return t.String() + " " + as[0]
}

func (t Table) String() string {
	return string(t)
}

type junction struct {
	left  Table
	right Table
}

func (j junction) swap() junction {
	return junction{left: j.right, right: j.left}
}

var juncs = map[junction]Table{
	{
		left:  TblWarehouse,
		right: TblPopularity,
	}: TblJuncWarehousePopularity,
	{
		left:  TblWarehouse,
		right: TblCommodity,
	}: TblJuncWarehouseCommodity,
	{
		left:  TblWarehouse,
		right: TblEquipment,
	}: TblJuncWarehouseEquipment,
	{
		left:  TblWarehouse,
		right: TblSystemParam,
	}: TblJuncWarehouseSystemParam,
	{
		left:  TblCell,
		right: TblPopularity,
	}: TblJuncCellPopularity,
	{
		left:  TblCell,
		right: TblCommodity,
	}: TblJuncCellCommodity,
	{
		left:  TblCell,
		right: TblEquipment,
	}: TblJuncCellEquipment,
}

// Junc gets the junction table between two parent tables.
func (t Table) Junc(right Table) (junc Table, ok bool) {
	j := junction{t, right}
	if junc, ok = juncs[j]; !ok {
		junc, ok = juncs[j.swap()]
	}

	return junc, ok
}

func (t Table) IsJunc() bool {
	switch t { //nolint:exhaustive //there is a default val for that
	case TblJuncWarehousePopularity, TblJuncWarehouseEquipment, TblJuncWarehouseCommodity, TblJuncWarehouseSystemParam:
		return true
	default:
		return false
	}
}

func (t Table) IDAsForKey() string {
	if t.IsJunc() {
		return ""
	}

	return fmt.Sprintf("%s_id", t.String())
}

type SQLAction string

const (
	Insert        SQLAction = "insert"
	Update        SQLAction = "update"
	PartialUpdate SQLAction = "patch"
	Delete        SQLAction = "delete"
	Replace       SQLAction = "replace"
	Select        SQLAction = "select"
)

func (a SQLAction) String() string {
	return string(a)
}
