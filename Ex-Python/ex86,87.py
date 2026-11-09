p=[]
i=[]

m=[[0,0,0],[0,0,0],[0,0,0]]
for l in range(0,3):
    for c in range(0,3):
        m[l][c]=int(input(f'Digite o valor {l+1},{c+1} da matriz :'))
print('___'*10)
for l in range(0,3):
    for c in range(0,3):
        print(f'{m[l][c]:^10}', end='')
        if m[l][c]%2==0:
            p.append(m[l][c])
        if m[l][c]%2==1:
            i.append(m[l][c])
    print()
print(f'A soma dos valores pares dessa matriz é {sum(p)}')
print(f'A soma dos valores da terceira coluna dessa matriz é {sum(m[2])}')
print(f'O maior valor da segunda linha dessa matriz é {max(m[1])}')
print('-END-'*10)