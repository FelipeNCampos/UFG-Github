r = 'S'
i=[]
b = []
p = []
while r =='S':
    n = int(input('Digite um numero a ser gravado : '))
    b.append(n)
    if n%2==0:
        p.append(n)
    if n%2!=0 :
        i.append(n)
    r=str(input('Deseja continuar (S/N) ?')).upper().strip()
print(f'Os numeros digitados foram {b}')
print(f'A lista de numeros pares foi {p}')
print(f'A lista de numeros impares foi {i}')