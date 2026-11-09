from ex108m import moedas

p = float(input(f'Digite o preço : R$'))
print(f'O dobro de R${p:.2f} é {moedas.dobro(p)}')
print(f'A metade de R${p:.2f} é {moedas.metade(p)}')
print(f'O valor de R${p:.2f} com aumento de 13% é {moedas.aumentar(p,13)}')
print(f'O valor de R${p:.2f} com redução de 7% é {moedas.diminuir(p,7)}')