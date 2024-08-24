def cpu_stress():
    counter = 0
    while True:
        # Realizar cálculos
        for _ in range(10**5):
            _ = 1 + 1
        
        # Pausa artificial
        counter += 1
        if counter % 33 == 0:  # Aproximadamente cada 3% de uso
            for _ in range(10**6):
                pass  # Esta pausa simula el sleep

if __name__ == "__main__":
    print("Iniciando estrés de CPU controlado")
    cpu_stress()