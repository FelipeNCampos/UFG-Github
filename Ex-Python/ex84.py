p=[]
d=[]
cpt = 0
pe= []
mpe = []
mle = []
while True:
    d.append(str(input('Digite o nome da primeira pessoa a ser cadastrada : ')).strip())
    d.append(float(input('Digite o peso da pessoa a ser cadastrada : ')))
    pe.append(d[1])
    p.append(d[:])
    d.clear()
    cpt += 1
    r=str(input('Deseja cadastrar mais uma pessoa (S/N)? ').strip().upper())
    if r=='N':
        break
print('==='*18)
print(f'{'Nome':^24}{'Peso':^24}')
print('___'*18)
for x, y in enumerate(p):
    print(f'{(y[0]).title():^24}', end='')
    print(f'{(y[1]):^24}')
print('___'*18)
print(f'Foram cadastradas um total de {cpt} pessoas')
for k in p:
    print (k[1])
    if k[1]==max(pe):
        mpe.append(k[0])
    if k[1]==min(pe):
        mle.append(k[0])
print(f'As pessoas mais pessadas registradas foram {mpe} com um total de {max(pe)} kilos.')
print(f'As pessoas mais leves registradas foram {mle} com um total de {min(pe)} kilos.')