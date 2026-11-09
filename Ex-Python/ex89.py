d = []
cad = []
while True:
    d.append(input('Digite o nome do aluno a ser cadastrado : '))
    for k in range(0,2):
        d.append(int(input(f'Digite a {k+1}* nota do aluno {d[0]} : ')))
    cad.append(d[:])
    d.clear()
    r = str(input('Deseja cadastrara mais um aluno (S/N) ?')).upper().strip()
    if r =='N':
        break
print('---'*10)
for c,x in enumerate(cad):
   print(f'As notas do aluno {x[0]} foram as seguintes :')
   for k in range(0,2):
        print(f'{c}*bimestre : {x[k]}')