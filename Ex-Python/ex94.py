pp={}
cad=[]
sum = 0
while True:
    pp['Nome']=input('Digite o nome da pessoa a ser cadastrada :').title()
    pp['Sexo']=input('Digite o sexo da pessoa (F/M):').strip().upper()
    while pp['Sexo']!='M' and pp['Sexo']!='F':
        print('Sexo inválido, digite novamente (M/F)')
    if pp['Sexo']=='M':
        pp['Idade']=int(input(f'Digite a idade do {pp['Nome']}:'))
    elif pp['Sexo']=='F':
        pp['Idade']=int(input(f'Digite a idade da {pp['Nome']}:'))
    cad.append(pp.copy())
    pp.clear()
    r=input('Deseja cadastrar mais uma pessoa (S/N) ?').strip().upper()
    if r=='N':
        break
print('-=-'*10)
print(f'Foram cadastradas um total de {len(cad)} pessoas.')
for c in cad:
    sum += c['Idade']
m = sum/int(len(cad))
print('As mulheres cadastradas foram :', end='')
for c in cad:
    if c['Sexo']=='F':
        print(f'{c['Nome']} ', end='')
print()
print('Os homens cadastrados foram :', end='')
for c in cad:
    if c['Sexo']=='M':
        print(f'{c['Nome']} ', end='')
print()
print(f'A idade media dos cadastrados foi de {m} anos, sendo os maiores que essa média:', end='')
for c in cad:
    if int(c['Idade'])>=m:
        print(f'{c['Nome']} ', end='')
print()