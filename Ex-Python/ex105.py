def notas(*n,Situação=False):
    sum = 0
    cad={}
    for c in n:
        sum += c
    #dic
    cad['Many_values']=len(n)
    cad['Maior nota']=max(n)
    cad['Media']=sum/len(n)
    if Situação==True:
        if cad['Media']<6:
            cad['Situação']='RUIM'
        if cad['Media']>=6 and cad['Media']<=7:
            cad['Situação']='RAZOÁVEL'
        if cad['Media']>7:
            cad['Situação']='BOM'
    return cad

r = notas(10,10,10,10,Situação=True)
print (r)