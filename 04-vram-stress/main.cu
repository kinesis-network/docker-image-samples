#include <cuda_runtime.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <time.h>

int main(int argc, char **argv) {
    int device_count = 0;
    cudaGetDeviceCount(&device_count);
    if (device_count == 0) {
        printf("No CUDA devices found.\n");
        return 1;
    }

    float min_ratio = 0.5f, max_ratio = 0.9f;
    int sleep_sec = 10;
    if (argc > 1) min_ratio = atof(argv[1]);
    if (argc > 2) max_ratio = atof(argv[2]);
    if (argc > 3) sleep_sec = atoi(argv[3]);

    srand(time(NULL));
    printf("Detected %d GPU(s)\n", device_count);
    printf("VRAM stress range: %.0f%%–%.0f%% (cycle every %d sec)\n",
           min_ratio * 100, max_ratio * 100, sleep_sec);

    while (1) {
        for (int i = 0; i < device_count; i++) {
            cudaSetDevice(i);
            size_t free_bytes = 0, total_bytes = 0;
            cudaMemGetInfo(&free_bytes, &total_bytes);

            float ratio = min_ratio + (rand() / (float)RAND_MAX) * (max_ratio - min_ratio);
            size_t alloc_bytes = (size_t)(free_bytes * ratio);

            printf("[GPU %d] Allocating %.2f GB (%.0f%% of free VRAM)\n",
                   i, alloc_bytes / (1024.0 * 1024 * 1024), ratio * 100);

            void *d_ptr = nullptr;
            cudaError_t err = cudaMalloc(&d_ptr, alloc_bytes);
            if (err != cudaSuccess) {
                printf("[GPU %d] Allocation failed: %s\n", i, cudaGetErrorString(err));
                continue;
            }

            cudaMemset(d_ptr, 0, alloc_bytes);
            printf("[GPU %d] Load active.\n", i);
        }

        printf("Sleeping for %d sec...\n", sleep_sec);
        sleep(sleep_sec);

        for (int i = 0; i < device_count; i++) {
            cudaSetDevice(i);
            cudaDeviceReset();
        }

        printf("Cycle complete. Restarting load...\n");
    }

    return 0;
}
