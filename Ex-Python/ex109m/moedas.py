def aumentar(p,v,form=True):
    if form==False:
        return p*((100+v)/100)
    if form==True:
        f = '{:.2f}'.format(p*((100+v)/100))
        return f


def diminuir(p,v,form=True):
    if form==False:
        return p*((100-v)/100)
    if form==True:
        f = '{:.2f}'.format(p*((100-v)/100))
        return f

def dobro (p,form=True):
    if form==False:
        return p*2
    if form==True:
        f = '{:.2f}'.format(p*2)
        return f


def metade(p,form=True):
    if form==False:
        return p/2
    if form==True:
        f = '{:.2f}'.format(p/2)
        return f

