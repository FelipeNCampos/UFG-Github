def leiaInt(n):
    while not n.isnumeric():
        n=(input(('\033[031m ERRO, digite um valor inteiro válido : \033[m')))
    else:
        return n


n = leiaInt(input('Digite um numero inteiro:'))
print(f'Você digitou o numero inteiro {n}')