from ex107m import moedas

p = float(input(f'Digite o preço : R$'))
print(f'O dobro de R${p} é {moedas.dobro(p)}')
print(f'A metade de R${p} é {moedas.metade(p)}')
print(f'O valor de R${p} com aumento de 13% é {moedas.aumentar(p,13)}')
print(f'O valor de R${p} com redução de 7% é {moedas.diminuir(p,7)}')