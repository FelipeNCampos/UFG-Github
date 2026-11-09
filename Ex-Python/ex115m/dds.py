def arqexiste(nome):
    try:
        a = open(nome,'rt')
        a.close()
    except FileNotFoundError:
        return False
    else:
        return True
def criararq(nome):
    try:
        a = open(nome, 'w+')
        a.close()
    except:
        print('\033[031mERRO\033[m ao criar o arquivo.')