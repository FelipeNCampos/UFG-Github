def fac(n,show):
    fac=1
    for c in range(1,n+1):
        fac *=c
        if show == True:
            print(f'{c} ', end='')
            if c<n:
                print(f'x ', end='')
            if c==n:
                print('= ', end='')
    print(fac)

print('---'*10)
n=int(input('Digite o numero a ser fatorado :'))
resp = input('Deseja mostrar o processo (S/N) ? ').strip().upper()
if resp=='S':
    show=True
else :
    show=False
fac(n,show)
