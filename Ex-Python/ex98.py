def contador (x,y,z):
    if z <0:
        z*=-1
    if z ==0:
        z=1
    if x>=y:
        print(x)
        while x>y:
            x-=z
            print(x)
    if y>x:
        print(x)
        while y>x:
            x+=z
            print(x)


x=int(input('Digite o numero que deseja começar a contagem :'))
y=int(input('Digite o numero que deseja terminar a contagem :'))
z=int(input('Digite o de quantos em quantos passos deseja :'))
contador(x,y,z)