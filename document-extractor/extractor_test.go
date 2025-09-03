package documentextractor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var pdfText = `H.01
CUENTA DE AHORROS
4884 1970 4694
INFORME DEL MES:MARZO /2021
Apreciado Cliente
JUAN CARLOS BALLESTEROS ROMERO
juanballesteros2001@gmail.com
??!!????##????$$????KK????&&????==????@@????++????WW????::??
Saldo Anterior $0.00
M�s Cr�ditos $947,332.64
Menos D�bitos $100.00
Nuevo Saldo $947,232.64
Saldo Promedio $207,582.90
Saldo Total Bolsillo $0.00
EXTRACTO CUENTA DE AHORROS
Fecha Valor Doc. Clase de Movimiento Oficina
10 03 $ 227,132.00+ 2613 Abono En Cuenta Por Pago De Nomina. PORTAL-EMPRESARIAL
10 03 $ 100.00- 6308 Compra BANCOLOMBIA Compras y Pagos PSE
30 03 $ 720,183.00+ 6789 Abono En Cuenta Por Pago De Nomina. PORTAL-EMPRESARIAL
31 03 $ 17.64+ 0000 Rendimientos Financieros.
Este producto cuenta con seguro de dep�sitos
Cualquier diferencia con el saldo, favor comunicarla a nuestra revisor�a fiscal KPMG Ltda. A.A. 77859 de Bogot�.
Recuerde que usted tambi�n cuenta con nuestro Defensor del Consumidor Financiero: Carlos Mario Serna Direcci�n: Calle 72 No. 6 - 30 Piso 18 en Bogot�. PBX: 6092013 Fax: 4829715 Correo Electr�nico:
defensordelcliente@davivienda.com
Para mayor informaci�n en www.davivienda.com
??abcdefghijklmnop ??
Banco Davivienda S.A NIT.860.034.313-7`

func TestDaviviendaExtractor(t *testing.T) {
	c := require.New(t)

	prevCurentExtractor := currentTextExtractor
	currentTextExtractor = func(pathPDF, password string) string { return pdfText }

	defer func() {
		currentTextExtractor = prevCurentExtractor
	}()

	extractor := &DaviviendaExtract{
		Password: "1005944768",
	}

	movements := extractor.GetMovements("dummypath")

	c.Len(movements, 4)
}
