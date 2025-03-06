#include <csignal>
#include <cstdio>
#include <chrono>
#include <thread>

const auto kSleepInterval = std::chrono::seconds(1);
bool exit_program = false;

int main(int argc, char *argv[]) {

  std::signal(SIGINT, [](int) {
    printf("Received SIGINT\n");
    exit_program = true;
  });

  while (!exit_program) {
    printf("Hello!\n");
    std::this_thread::sleep_for(kSleepInterval);
  }
}
