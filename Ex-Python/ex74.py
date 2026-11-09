from random import randint 
t = (randint(0,999),randint(0,999),randint(0,999),randint(0,999),randint(0,999))
print (t)
ma = t[0]
me = t[0]
for c in range(0,4):
    if t[c]<t[c+1]:
        ma = t[c+1]
    if t[c]>t[c+1]:
        me = t[c+1]
print('Maior {}'.format(ma))
print('Menor {}'.format(me))

