import urllib
import urllib.request

try :
    site = urllib.request.urlopen('https://www.youtube.com/')
except:
    print('\033[031m ERRO \033[m')
else :
    print('\033[032m SUCESSO \033[m')
    print(site.read())