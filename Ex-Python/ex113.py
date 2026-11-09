def leiaint(n1):
    while True:
        try:
            n1 = int(input('Digite um numero inteiro : '))
            print(f'Você digitou o numero inteiro {n1}')
            return n1
        except:
            print(f'Você digitou um valor inteiro \033[031minválido\033[m, tente novamente!')

def leiafloat(n2):
    while True:
        try:
            n2 = float(input('Digite um numero real : '))
            print(f'Você digitou o numero real {n2:.2f}')
            return n2
        except:
            print(f'Você digitou um valor real\033[031m inválido \033[m, tente novamente!')


n1 = n2 = 0
leiaint(n1)
leiafloat(n2)