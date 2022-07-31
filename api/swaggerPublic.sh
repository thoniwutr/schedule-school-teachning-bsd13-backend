cp swagger.yaml swaggerPublic.yaml

perl -pi -w -e 's/dto.//g;' swaggerPublic.yaml
perl -pi -w -e 's/util.//g;' swaggerPublic.yaml
