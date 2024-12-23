package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

type PurchaseRecord struct {
	NamaCustomer string
	NamaVendor   string
	HargaBeli    int
	HargaJual    int
	Fee          int
}

func createPurchaseData(data [][]string) []PurchaseRecord {
	var purchaseList []PurchaseRecord
	for i, line := range data {
		if i > 0 {
			var rec PurchaseRecord
			for j, field := range line {
				switch j {
				case 3:
					rec.NamaCustomer = field
				case 4:
					rec.NamaVendor = field
				case 6:
					v, _ := strconv.Atoi(field)
					rec.HargaBeli = v
				case 7:
					v, _ := strconv.Atoi(field)
					rec.HargaJual = v
				case 8:
					v, _ := strconv.Atoi(field)
					rec.Fee = v
				}
			}
			purchaseList = append(purchaseList, rec)
		}
	}

	return purchaseList
}

const (
	ARTPAYPPOB            = "ArtPay PPOB"
	AGRATEKPPOB           = "Agratek PPOB"
	KIOSDESAPPOB          = "Kios Desa"
	TOKOPEDIAPPOB         = "Tokopedia PPOB"
	TOKOPEDIADISBURSEMENT = "Tokopedia Disbursement"
	IDS                   = "Inovasi Daya Solusi"
)

const (
	LINKAJA    = "Link Aja"
	RAJABILLER = "RAJABILLER"
	OPENSIPKD  = "OpenSIPKD"
	NICEPAY    = "NICE PAY"
	IDSVENDOR  = "IDS Vendor"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Masukkan file CSV cuk!")
		os.Exit(1)
	}

	args := os.Args
	argCsv := args[1]
	if !strings.Contains(argCsv, ".csv") {
		fmt.Println("File harus CSV cuk!")
		os.Exit(1)
	}

	f, err := os.Open(argCsv)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err.Error())
	}

	purchaseList := createPurchaseData(data)

	var (
		artpay                = 0
		kiosdesa              = 0
		tokopediaPPOB         = 0
		tokopediaDisbursement = 0
		agratekPPOB           = 0
		inovasiDayaSolusi     = 0
		totalJual             = 0
	)

	var (
		artpayCount                = 0
		kiosdesaCount              = 0
		tokopediaPPOBCount         = 0
		tokopediaDisbursementCount = 0
		agratekPPOBCount           = 0
		inovasiDayaSolusiCount     = 0
		totalTransaksiCount        = 0
	)

	var (
		linkaja    = 0
		rajabiller = 0
		opensipkd  = 0
		nicepay    = 0
		idsVendor  = 0
		totalBeli  = 0
	)

	var (
		linkajaCount    = 0
		rajabillerCount = 0
		opensipkdCount  = 0
		idsVendorCount  = 0
		nicepayCount    = 0
	)

	for _, v := range purchaseList {
		totalBeli += v.HargaBeli
		totalJual += v.HargaJual
		totalTransaksiCount += 1

		// FOR CALCULATE CUSTOMER
		switch v.NamaCustomer {
		case ARTPAYPPOB:
			artpay += v.HargaJual
			artpayCount += 1
		case KIOSDESAPPOB:
			kiosdesa += v.HargaJual
			kiosdesaCount += 1
		case AGRATEKPPOB:
			agratekPPOB += v.HargaJual
			agratekPPOBCount += 1
		case TOKOPEDIAPPOB:
			tokopediaPPOB += v.HargaJual
			tokopediaPPOBCount += 1
		case TOKOPEDIADISBURSEMENT:
			tokopediaDisbursement += v.HargaJual
			tokopediaDisbursementCount += 1
		case IDS:
			inovasiDayaSolusi += v.HargaJual
			inovasiDayaSolusiCount += 1

		}

		// FOR CALCULATE VENDOR
		switch v.NamaVendor {
		case LINKAJA:
			linkaja += v.HargaBeli
			linkajaCount += 1
		case RAJABILLER:
			rajabiller += v.HargaBeli
			rajabillerCount += 1
		case OPENSIPKD:
			opensipkd += v.HargaBeli
			opensipkdCount += 1
		case NICEPAY:
			nicepay += v.HargaBeli
			nicepayCount += 1
		case IDSVENDOR:
			idsVendor += v.HargaBeli
			idsVendorCount += 1
		}
	}

	customerTotal := artpay + kiosdesa + tokopediaPPOB + tokopediaDisbursement + inovasiDayaSolusi + agratekPPOB
	vendorTotal := linkaja + rajabiller + opensipkd + nicepay + idsVendor
	feeTotal := totalJual - totalBeli
	customerTotalCount := inovasiDayaSolusiCount + artpayCount + kiosdesaCount + agratekPPOBCount + tokopediaPPOBCount + tokopediaDisbursementCount
	vendorTotalCount := linkajaCount + rajabillerCount + opensipkdCount + nicepayCount + idsVendorCount

	// BALANCE CHECK
	var isVendorBalance bool = totalBeli == vendorTotal
	var isVendorBalanceNaration string
	var isMerchantBalance bool = totalJual == customerTotal
	var isMerchantBalanceNaration string

	if isVendorBalance {
		isVendorBalanceNaration = "BALANCE!"
	} else {
		isVendorBalanceNaration = "NOT BALANCE!"
	}

	if isMerchantBalance {
		isMerchantBalanceNaration = "BALANCE!"
	} else {
		isMerchantBalanceNaration = "NOT BALANCE!"
	}

	// TABLE FOR VENDOR
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendRows([]table.Row{
		{"Total Transaksi", totalTransaksiCount},
		{"Total Beli", totalBeli},
		{"Total Jual", totalJual},
		{"Total Fee", feeTotal},
	})
	t.SetCaption("REKAP TRANSACTION\n")
	t.Render()

	t.ResetHeaders()
	t.ResetRows()
	t.ResetFooters()

	t.AppendHeader(table.Row{"#", "Vendor", "Total Transaksi", "Total Amount"})
	t.AppendRows([]table.Row{
		{1, "Link Aja", linkajaCount, linkaja},
		{2, "Rajabiller", rajabillerCount, rajabiller},
		{3, "OpenSIPKD", opensipkdCount, opensipkd},
		{4, "Nice Pay", nicepayCount, nicepay},
		{5, "IDS Vendor", idsVendorCount, idsVendor},
	})
	t.AppendFooter(table.Row{"", "Total", vendorTotalCount, vendorTotal})
	t.SetCaption("VENDOR TRANSACTION\nMembandingkan Total Beli Dengan Total Vendor\nTotal Beli: %v\nTotal Vendor: %v\nHasil: %v\n", totalBeli, vendorTotal, isVendorBalanceNaration)
	t.Render()

	t.ResetHeaders()
	t.ResetRows()
	t.ResetFooters()

	// TABLE FOR MERCHANT
	t.AppendHeader(table.Row{"#", "Merchant", "Total Transaksi", "Total Amount"})
	t.AppendRows([]table.Row{
		{1, "IDS", inovasiDayaSolusiCount, inovasiDayaSolusi},
		{2, "Artpay", artpayCount, artpay},
		{3, "Kiosdesa", kiosdesaCount, kiosdesa},
		{4, "Agratek PPOB", agratekPPOBCount, agratekPPOB},
		{5, "Tokopedia PPOB", tokopediaPPOBCount, tokopediaPPOB},
		{6, "Tokopedia Disbursement", tokopediaDisbursementCount, tokopediaDisbursement},
	})
	t.AppendFooter(table.Row{"", "Total", customerTotalCount, customerTotal})
	t.SetCaption("MERCHANT TRANSACTION\nMembandingkan Total Jual Dengan Total Merchant\nTotal Jual: %v\nTotal Merchant: %v\nHasil: %v\n", totalJual, customerTotal, isMerchantBalanceNaration)
	t.Render()
}
