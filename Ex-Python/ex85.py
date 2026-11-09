i=[]
p=[]

for c in range(0,7):
    n = int(input('Digite um numero :'))
    if n%2==0:
        p.append(n)
    if n%2==1:
        i.append(n)
p.sort()
i.sort()
print(f'Os impares sao : {i}')
print (f'Os pares sao : {p}')