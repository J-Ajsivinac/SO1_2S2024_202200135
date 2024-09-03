import time

def cpu_stress():
    while True:
        # Realiza cálculos intensivos
        for _ in range(198000):
            _ = 3.1415 * 2.7182
        
        # Pausa para reducir el uso de CPU
        time.sleep(0.92)

if __name__ == "__main__":
    cpu_stress()
