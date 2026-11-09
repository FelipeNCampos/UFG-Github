n = str(input('Digite o nome da pessoa que ira sacar o dinheiro : '))
print('-=-'*20)
print ('Bem vindo ao caixa eletronico {}'.format(n))
print('-=-'*20)
notas = [50,20,10,1]
v = int(input(f'Digite o valor a ser sacado {n} : '))
for i in range(0,4):
    nmt = v//notas[i]
    if nmt != 0 :
        print('Sao {} cedulas de {}'.format(nmt,notas[i]))
    v -= nmt*notas[i]
