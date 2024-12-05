# 🚀 FastCopy: High-Performance File Copy Tool 🖹

**FastCopy** is a blazing-fast file copy tool written in **Go**, optimized for multi-core systems. Designed for speed, efficiency, and simplicity, it utilizes parallel processing and system-level APIs to deliver superior file transfer performance. Whether you're copying a small document or a massive ISO file, GoFastCopy gets the job done faster than ever! ⚡

---

## **✨ Features**
- 🏎️ **High-Speed Copying**: Leverages multi-threaded processing for unmatched performance.
- 📊 **Real-Time Progress**: Displays live progress updates with percentage and speed (MB/s).
- 💻 **Windows-Optimized**: Built using Windows system calls for direct, efficient I/O.
- 🛠️ **Customizable**: Adjust the number of threads to suit your system's capabilities.
- 🔗 **Chunk-Based Copy**: Efficiently handles large files by dividing them into manageable chunks.

---

## **📥 Installation**

### **1. Prerequisites**
- **Go** (version 1.18+): Download from [golang.org](https://golang.org/dl/).
- A **Windows** machine (FastCopy currently supports Windows only).

### **2. Clone the Repository**
```bash
git clone https://github.com/kapilgarg/fastcopy.git
cd FastCopy
```

### **3. Build the project**
```
go build -o fastcopy main.go
```

## **🚀 Usage**
```
fastcopy <source_file_path> <destination_directory> <number_of_workers>
```

## **📊 Output**
Copying... 25% | 120.45 MB/s  
Copying... 50% | 118.23 MB/s  
Copying... 75% | 115.67 MB/s  
Copying... 100% | 0.00 MB/s  
Copy complete!

## **⚙️ How It Works**
1. File Splitting: The file is divided into chunks based on the number of workers.
2. Parallel Copying: Each worker thread copies its assigned chunk using direct I/O operations.
3. Progress Tracking: Displays real-time percentage and transfer speed.

## **🛠️ Performance Tuning**
1. Adjust the bufferSize constant in the code (default: 4 MB) for optimal disk I/O.
2. Increase the number of workers for multi-core processors to improve speed.
3. Benchmark for your system to find the sweet spot for performance.

