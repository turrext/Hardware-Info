package main

import (
	"fmt"
	"image/color"
	"math"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/jaypipes/ghw"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type Information struct {
	Hostname   string
	Platform   string
	TotalDisks string
	DiskTotal  string
	RAM        uint64
	CPU        string
	DiskUsed   string
}

type GpuInfo struct {
	Name string
	VRam string
}

func main() {
	var TotalDisksA int
	myApp := app.New()
	myWindow := myApp.NewWindow("Computer Status")
	hostStat, _ := host.Info()
	cpuStat, _ := cpu.Info()
	vmStat, _ := mem.VirtualMemory()
	info := new(Information)

	block, err := ghw.Block()
	if err != nil {
		fmt.Printf("Error getting block storage info: %v", err)
	}
	for diskn, _ := range block.Disks {
		TotalDisksA = diskn + 1
	}
	info.TotalDisks = strconv.Itoa(TotalDisksA)
	info.DiskTotal = "69"
	info.DiskTotal = "60"
	info.Hostname = hostStat.Hostname
	info.Platform = hostStat.Platform
	info.CPU = cpuStat[0].ModelName
	info.RAM = vmStat.Total / 1024 / 1024

	AddBtns(myApp, myWindow, info)

	myWindow.Resize(fyne.NewSize(400, 600))
	myWindow.SetMaster()
	myWindow.CenterOnScreen()

	myWindow.ShowAndRun()

}

func AddBtns(myApp fyne.App, myWindow fyne.Window, info *Information) {

	Info := exec.Command("cmd", "/C", "wmic path win32_VideoController get name")
	Info.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	History, _ := Info.Output()
	replace := strings.Replace(string(History), "Name", "", -1)
	replace2 := strings.Replace(replace, "LuminonCore IDDCX Adapter", "", -1)
	Info2 := exec.Command("cmd", "/C", "wmic path Win32_videocontroller get adapterram")
	Info2.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	History2, _ := Info2.Output()
	replace3 := strings.Replace(string(History2), "AdapterRAM", "", -1)
	split3 := strings.Split(replace3, "\n")
	split2 := strings.Split(replace2, "\n")
	var gpuar []string
	gpuss := 0
	print(len(split2))
	for x, vram := range split3 {
		gpuname := strings.ReplaceAll(split2[x], " ", "")
		gpuvram := strings.ReplaceAll(vram, " ", "")
		if gpuname == "" || vram == "  " || vram == "" || len(gpuname) < 4 || len(vram) < 4 {

		} else {
			var gpu GpuInfo
			gpu.Name = gpuname
			gpu.VRam = gpuvram
			gpuar = append(gpuar, gpu.Name)
			gpuar = append(gpuar, gpu.VRam)
			gpuss = gpuss + 1
		}

	}
	CPU1 := widget.NewCard("CPU Info", info.CPU, widget.NewButton("More", func() {
		myCPU := myApp.NewWindow("CPU Info")
		openCPU(myApp, myWindow, myCPU, info)

	}))
	RamGB := float64(info.RAM) / float64(1024)
	RamGB = roundFloat(RamGB, 1)
	RamStr := strconv.FormatFloat(RamGB, 'f', -1, 64)
	RAM1 := widget.NewCard("Ram Info", "Memory: "+RamStr+" GB", widget.NewButton("More", func() {
		myRAM := myApp.NewWindow("Ram Info")
		openRAM(myApp, myWindow, myRAM, info)

	}))

	PC1 := widget.NewCard("PC Info", "Hostname: "+info.Hostname,
		widget.NewButton("There are: "+info.TotalDisks+" Storage Devices On the Machine!", func() {
			myDisk := myApp.NewWindow("Additional Info")
			openDISK(myApp, myWindow, myDisk, info)
		}))

	gpunum1 := strconv.Itoa(gpuss)
	GPU1 := widget.NewCard("GPU Info", "There are "+gpunum1+" Graphic Card(s)", widget.NewButton("Info", func() {
		myGPU := myApp.NewWindow("GPU Info")
		openGPU(myApp, myWindow, myGPU, gpuar)
	}))

	exit := widget.NewButton("Quit", func() {
		myWindow.Close()

	})
	myWindow.SetContent(container.NewBorder(container.New(layout.NewVBoxLayout(), RAM1, CPU1, GPU1, PC1), container.New(layout.NewVBoxLayout(), exit), nil, nil))
}

func openCPU(app fyne.App, myWindow fyne.Window, myCPU fyne.Window, info *Information) {
	var v []cpu.InfoStat
	var err error
	var cpucores int
	var Family string
	var PhysicalID string
	var CoreID string
	var ModelName string
	var Mhz float64
	var Microcode string

	if v, err = cpu.Info(); err != nil {
		panic(err)
		return
	}
	for i, cpu := range v {
		fmt.Println("Ran ", i, " Times")
		cpucores = i
		Family = cpu.Family
		PhysicalID = cpu.PhysicalID
		CoreID = cpu.CoreID
		ModelName = cpu.ModelName
		Mhz = cpu.Mhz
		Microcode = cpu.Microcode
		break
	}

	exitcpu := widget.NewButton("Exit", func() {

		myCPU.Hide()

	})
	cpuStat, _ := cpu.Info()
	strcpu := strconv.Itoa(cpucores)
	strmhz := strconv.FormatFloat(Mhz, 'f', -1, 64)
	mhztoint, _ := strconv.Atoi(strmhz)
	fmt.Println(mhztoint)
	ghzint := float64(mhztoint) / float64(1000)
	fmt.Println(ghzint)
	ghzstr := strconv.FormatFloat(ghzint, 'f', -1, 64)
	text1 := canvas.NewText("CPU: "+strcpu, color.Black)
	text2 := canvas.NewText("Family: "+Family, color.Black)
	text3 := canvas.NewText("Model: "+cpuStat[0].Model, color.Black)
	text5 := canvas.NewText("Physical ID: "+PhysicalID, color.Black)
	text6 := canvas.NewText("CoreID: "+CoreID, color.Black)
	text8 := canvas.NewText("Modelname: "+ModelName, color.Black)
	text9 := canvas.NewText("GHZ: "+ghzstr, color.Black)
	text12 := canvas.NewText("Microcode: "+Microcode, color.Black)

	myCPU.SetContent(container.NewBorder(container.New(layout.NewVBoxLayout(), text8, text1, text2, text3, text5, text6, text9, text12), container.New(layout.NewVBoxLayout(), exitcpu), nil, nil))
	myCPU.Resize(fyne.NewSize(400, 400))
	myCPU.Show()
}

func openRAM(app fyne.App, myWindow fyne.Window, myRAM fyne.Window, info *Information) {
	v, _ := mem.VirtualMemory()

	exitram := widget.NewButton("Exit", func() {

		myRAM.Hide()

	})
	intramsGB := float64(v.Total) / float64(1024) / float64(1024) / float64(1024)
	intrams := strconv.FormatFloat(intramsGB, 'f', 1, 64)
	availRamGB := float64(v.Available) / float64(1024) / float64(1024) / float64(1024)
	intavailable := strconv.FormatFloat(availRamGB, 'f', 1, 64)
	usedRamGB := float64(v.Used) / float64(1024) / float64(1024) / float64(1024)
	intused := strconv.FormatFloat(usedRamGB, 'f', 1, 64)
	usedRamPercent := strconv.Itoa(int(v.UsedPercent))
	text1 := canvas.NewText("Total: "+intrams+" GB", color.Black)
	text2 := canvas.NewText("Available Memory: "+intavailable+" GB", color.Black)
	text3 := canvas.NewText("Used Memory: "+intused+" GB", color.Black)
	text5 := canvas.NewText("Used Percent: "+usedRamPercent+"%", color.Black)

	myRAM.SetContent(container.NewBorder(container.New(layout.NewVBoxLayout(), text1, text2, text3, text5), container.New(layout.NewVBoxLayout(), exitram), nil, nil))
	myRAM.Resize(fyne.NewSize(400, 400))
	myRAM.Show()
}

func openDISK(app fyne.App, myWindow fyne.Window, myDISK fyne.Window, info *Information) {

	exitdisk := widget.NewButton("Exit", func() {

		myDISK.Hide()

	})
	block, err := ghw.Block()
	if err != nil {
		fmt.Printf("Error getting block storage info: %v", err)
	}
	//var storagedisks []string
	for i, disk := range block.Disks {
		fmt.Println("Type: ", i, ": ", disk.DriveType)
		fmt.Println("Model: ", i, ": "+disk.Model)
		fmt.Println("Vendor: ", i, ": "+disk.Vendor)
		fmt.Println("SizeBytes: ", i, ": ", strconv.FormatFloat((float64(disk.SizeBytes)/float64(1024)/float64(1024)/float64(1024)), 'f', 1, 64)+" GB")
	}
	text1 := canvas.NewText("Total:  GB", color.Black)
	myDISK.SetContent(container.NewBorder(container.New(layout.NewVBoxLayout(), text1), container.New(layout.NewVBoxLayout(), exitdisk), nil, nil))
	myDISK.Resize(fyne.NewSize(400, 400))
	myDISK.Show()
}

func openGPU(app fyne.App, myWindow fyne.Window, myGPU fyne.Window, GpuInfo []string) {

	exitgpu := widget.NewButton("Exit", func() {

		myGPU.Hide()

	})
	var gpu1 string
	var gpu2 string
	var gpu3 string
	var gpu4 string
	var x int
	for i, _ := range GpuInfo {

		if i%2 == 0 || i == 0 {
			if i == 0 {
				x = i + 1
				returneds := strings.TrimSpace(GpuInfo[x])
				n, _ := strconv.ParseInt(returneds, 10, 64)
				gpuvrammbf := float64(n) / float64(1024) / float64(1024)
				gpuvrammbs := strconv.FormatFloat(gpuvrammbf, 'f', 1, 64)
				gpu1 = GpuInfo[i] + " : " + gpuvrammbs + " MB"
			}
			if i == 2 {
				x = i + 1
				returneds := strings.TrimSpace(GpuInfo[x])
				n, _ := strconv.ParseInt(returneds, 10, 64)
				gpuvrammbf := float64(n) / float64(1024) / float64(1024)
				gpuvrammbs := strconv.FormatFloat(gpuvrammbf, 'f', 1, 64)
				gpu2 = GpuInfo[i] + " : " + gpuvrammbs + " MB"
			}
			if i == 4 {
				x = i + 1
				returneds := strings.TrimSpace(GpuInfo[x])
				n, _ := strconv.ParseInt(returneds, 10, 64)
				gpuvrammbf := float64(n) / float64(1024) / float64(1024)
				gpuvrammbs := strconv.FormatFloat(gpuvrammbf, 'f', 1, 64)
				gpu3 = GpuInfo[i] + " : " + gpuvrammbs + " MB"
			}
			if i == 6 {
				x = i + 1
				returneds := strings.TrimSpace(GpuInfo[x])
				n, _ := strconv.ParseInt(returneds, 10, 64)
				gpuvrammbf := float64(n) / float64(1024) / float64(1024)
				gpuvrammbs := strconv.FormatFloat(gpuvrammbf, 'f', 1, 64)
				gpu4 = GpuInfo[i] + " : " + gpuvrammbs + " MB"
			}

		}
	}
	if gpu1 != "" && gpu2 != "" && gpu3 != "" && gpu4 != "" {

		text1 := canvas.NewText("GPU: "+gpu1, color.Black)
		text2 := canvas.NewText("GPU: "+gpu2, color.Black)
		text3 := canvas.NewText("GPU: "+gpu3, color.Black)
		text4 := canvas.NewText("GPU: "+gpu4, color.Black)
		myGPU.SetContent(container.NewBorder(container.New(layout.NewVBoxLayout(), text1, text2, text3, text4), container.New(layout.NewVBoxLayout(), exitgpu), nil, nil))

	} else if gpu1 != "" && gpu2 != "" && gpu3 != "" {
		text1 := canvas.NewText("GPU: "+gpu1, color.Black)
		text2 := canvas.NewText("GPU: "+gpu2, color.Black)
		text3 := canvas.NewText("GPU: "+gpu3, color.Black)
		myGPU.SetContent(container.NewBorder(container.New(layout.NewVBoxLayout(), text1, text2, text3), container.New(layout.NewVBoxLayout(), exitgpu), nil, nil))

	} else if gpu1 != "" && gpu2 != "" {
		text1 := canvas.NewText("GPU: "+gpu1, color.Black)
		text2 := canvas.NewText("GPU: "+gpu2, color.Black)
		myGPU.SetContent(container.NewBorder(container.New(layout.NewVBoxLayout(), text1, text2), container.New(layout.NewVBoxLayout(), exitgpu), nil, nil))

	} else if gpu1 != "" {
		text1 := canvas.NewText("GPU: "+gpu1, color.Black)
		myGPU.SetContent(container.NewBorder(container.New(layout.NewVBoxLayout(), text1), container.New(layout.NewVBoxLayout(), exitgpu), nil, nil))

	} else {
		text1 := canvas.NewText("Unable to Fetch GPU Info", color.Black)
		myGPU.SetContent(container.NewBorder(container.New(layout.NewVBoxLayout(), text1), container.New(layout.NewVBoxLayout(), exitgpu), nil, nil))

	}
	print(gpu1 + gpu2 + gpu3 + gpu4)

	myGPU.Resize(fyne.NewSize(400, 400))
	myGPU.Show()
}

func String(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
