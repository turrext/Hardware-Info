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

type raminfo1 struct {
	TotalMemory       string
	AvailableMemory   string
	UsedMemory        string
	UsedMemoryPercent string
}

type cpu1 struct {
	Family     string
	PhysicalID string
	ModelName  string
	Mhz        float64
}

type Information struct {
	Hostname   string
	Platform   string
	TotalDisks string
	RAM        uint64
	CPU        string
}

type GpuInfo struct {
	Name string
	VRam string
}

////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////
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
	gpuar := getgpus()
	RamStr := ramgbstr(info.RAM)

	CPU1 := widget.NewCard("CPU", info.CPU, widget.NewButton("More", func() {
		myCPU := myApp.NewWindow("CPU Info")
		openCPU(myApp, myWindow, myCPU, info)

	}))
	RAM1 := widget.NewCard("Ram", "Memory: "+RamStr+" GB", widget.NewButton("More", func() {
		myRAM := myApp.NewWindow("Ram Info")
		openRAM(myApp, myWindow, myRAM, info)

	}))
	PC1 := widget.NewCard("PC", "Hostname: "+info.Hostname,
		widget.NewButton("There are: "+info.TotalDisks+" Storage Devices On the Machine!", func() {
			myDisk := myApp.NewWindow("Additional Info")
			openDISK(myApp, myWindow, myDisk, info)
		}))
	GPU1 := widget.NewCard("GPU", "Video Card(s) Information", widget.NewButton("View", func() {
		myGPU := myApp.NewWindow("GPU Info")
		openGPU(myApp, myWindow, myGPU, gpuar)
	}))
	exit := widget.NewButton("Quit", func() {
		myWindow.Close()

	})

	myWindow.SetContent(container.NewBorder(container.New(layout.NewVBoxLayout(), RAM1, CPU1, GPU1, PC1), container.New(layout.NewVBoxLayout(), exit), nil, nil))
}

////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////
func openCPU(app fyne.App, myWindow fyne.Window, myCPU fyne.Window, info *Information) {

	//////
	cpu := getcpu()

	////////
	ghzstr := FtoaGHZ(cpu.Mhz)
	text2 := canvas.NewText("Family: "+cpu.Family, color.Black)
	text5 := canvas.NewText("Physical ID: "+cpu.PhysicalID, color.Black)
	text8 := canvas.NewText("Modelname: "+cpu.ModelName, color.Black)
	text9 := canvas.NewText("GHZ: "+ghzstr, color.Black)

	////////
	exitcpu := widget.NewButton("Exit", func() {

		myCPU.Hide()

	})

	myCPU.SetContent(container.NewBorder(container.New(layout.NewVBoxLayout(), text8, text2, text5, text9), container.New(layout.NewVBoxLayout(), exitcpu), nil, nil))
	myCPU.Resize(fyne.NewSize(400, 400))
	myCPU.Show()
}
func openRAM(app fyne.App, myWindow fyne.Window, myRAM fyne.Window, info *Information) {

	exitram := widget.NewButton("Exit", func() {

		myRAM.Hide()

	})

	memory := getraminfo()

	text1 := canvas.NewText("Total: "+memory.TotalMemory+" GB", color.Black)
	text2 := canvas.NewText("Available Memory: "+memory.AvailableMemory+" GB", color.Black)
	text3 := canvas.NewText("Used Memory: "+memory.UsedMemory+" GB", color.Black)
	text5 := canvas.NewText("Used Percent: "+memory.UsedMemoryPercent+"%", color.Black)

	myRAM.SetContent(container.NewBorder(container.New(layout.NewVBoxLayout(), text1, text2, text3, text5), container.New(layout.NewVBoxLayout(), exitram), nil, nil))
	myRAM.Resize(fyne.NewSize(400, 400))
	myRAM.Show()
}
func openDISK(app fyne.App, myWindow fyne.Window, myDISK fyne.Window, info *Information) {
	rows := getdiskrows()
	exitdisk := widget.NewButton("Exit", func() {

		myDISK.Hide()

	})

	myDISK.SetContent(container.NewBorder(container.New(layout.NewVBoxLayout(), rows...), container.New(layout.NewVBoxLayout(), exitdisk), nil, nil))
	myDISK.Resize(fyne.NewSize(400, 400))
	myDISK.Show()
}
func openGPU(app fyne.App, myWindow fyne.Window, myGPU fyne.Window, GpuInfo []string) {

	gpucontent(GpuInfo, myGPU)

	myGPU.Resize(fyne.NewSize(400, 400))
	myGPU.Show()
}

////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////
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
func FtoaGHZ(var1 float64) string {

	ghzint := float64(var1) / float64(1000)
	ghzstr := strconv.FormatFloat(ghzint, 'f', -1, 64)
	return ghzstr

}
func getcpu() cpu1 {
	var d []cpu.InfoStat
	var err error
	var cpu1 cpu1
	if d, err = cpu.Info(); err != nil {
		fmt.Printf("%D\n", err)
		return cpu1
	}
	for _, cpu := range d {
		cpu1.Family = cpu.Family
		cpu1.PhysicalID = cpu.PhysicalID
		cpu1.ModelName = cpu.ModelName
		cpu1.Mhz = cpu.Mhz
		break
	}
	return cpu1
}
func gpucontent(GpuInfo []string, myGPU fyne.Window) {

	exitgpu := widget.NewButton("Exit", func() {

		myGPU.Hide()

	})

	var gpu1 string
	var gpu2 string
	var gpu3 string
	var gpu4 string
	var x int
	for i := range GpuInfo {

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
}
func getgpusname() string {

	Info := exec.Command("cmd", "/C", "wmic path win32_VideoController get name")
	Info.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	History, _ := Info.Output()
	replace := strings.Replace(string(History), "Name", "", -1)
	replace2 := strings.Replace(replace, "LuminonCore IDDCX Adapter", "", -1)
	return replace2
}
func getgpusvram() string {
	Info2 := exec.Command("cmd", "/C", "wmic path Win32_videocontroller get adapterram")
	Info2.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	History2, _ := Info2.Output()
	replace3 := strings.Replace(string(History2), "AdapterRAM", "", -1)
	return replace3
}
func getgpus() []string {
	names := getgpusname()
	vram := getgpusvram()
	split3 := strings.Split(vram, "\n")
	split2 := strings.Split(names, "\n")
	var gpuar []string
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
		}

	}
	return gpuar
}
func ramgbstr(ram uint64) string {
	RamGB := float64(ram) / float64(1024)
	RamGB = roundFloat(RamGB, 1)
	RamStr := strconv.FormatFloat(RamGB, 'f', -1, 64)
	return RamStr
}
func getraminfo() raminfo1 {
	var memory raminfo1
	v, _ := mem.VirtualMemory()
	intramsGB := float64(v.Total) / float64(1024) / float64(1024) / float64(1024)
	memory.TotalMemory = strconv.FormatFloat(intramsGB, 'f', 1, 64)
	availRamGB := float64(v.Available) / float64(1024) / float64(1024) / float64(1024)
	memory.AvailableMemory = strconv.FormatFloat(availRamGB, 'f', 1, 64)
	usedRamGB := float64(v.Used) / float64(1024) / float64(1024) / float64(1024)
	memory.UsedMemory = strconv.FormatFloat(usedRamGB, 'f', 1, 64)
	memory.UsedMemoryPercent = strconv.Itoa(int(v.UsedPercent))
	return memory
}
func getdiskrows() []fyne.CanvasObject {

	rows := []fyne.CanvasObject{}
	block, err := ghw.Block()
	if err != nil {
		fmt.Printf("Error getting block storage info: %v", err)
	}
	for _, disk := range block.Disks {

		disktype := fmt.Sprintf(disk.Model)
		rows = append(
			rows,
			container.New(
				layout.NewVBoxLayout(),
				widget.NewLabel("Disk Model: "+disktype),
				widget.NewLabel("Size: "+strconv.FormatFloat((float64(disk.SizeBytes)/float64(1024)/float64(1024)/float64(1024)), 'f', 1, 64)+" GB"),
				widget.NewLabel("Disk Vendor: "+disk.Vendor),
				widget.NewLabel("Drive Type: "+disk.DriveType.String())),
		)

	}
	return rows
}
